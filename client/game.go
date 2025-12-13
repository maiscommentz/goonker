package main

import (
	"Goonker/client/audio"
	"Goonker/client/ui"
	"Goonker/common"

	"encoding/json"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// States of the application
	sInit = iota
	sMainMenu
	sPlayMenu
	sWaitingGame
	sGamePlaying
	sGameWin
	sGameLose
	sGameDraw

	// Network configuration
	serverAddress = "ws://localhost:8080/ws" // goonker.saikoon.ch
	//roomId        = "87DY68"
	isBotGame = true
)

type Game struct {
	menu         *ui.MainMenu
	playMenu     *ui.PlayMenu
	waitingMenu  *ui.WaitingMenu
	state        int
	netClient    *NetworkClient
	grid         *ui.Grid
	audioManager *audio.AudioManager

	mySymbol common.PlayerID // 1 for X, 2 for O
	isMyTurn bool
}

// Init the game
func (g *Game) Init() {
	// Initialize network client
	g.netClient = NewNetworkClient()

	// Initialize the UI
	ui.Init()

	// Initialize UI elements
	g.initUIElements()

	// Initialize Audio Manager and sounds
	g.initAudio()

	// Play main menu music
	err := g.audioManager.LoadMusic("main_menu_music", "client/assets/main_menu.mp3")
	if err != nil {
		log.Println("Could not load music:", err)
	}
	g.audioManager.Play("main_menu_music")

	// Set initial state
	g.state = sMainMenu
}

// Update the game state every tick.
// Typically called every tick (1/60[s] by default).
func (g *Game) Update() error {

	// Always poll the network for incoming messages first
	g.handleNetwork()

	switch g.state {
	case sInit:
		g.Init()
	case sMainMenu:
		if g.menu.BtnPlay.IsClicked() {
			g.audioManager.Play("click_button")
			// Try to connect to server (Async)
			// Note: For WASM/Localhost testing use ws://localhost:8080/ws?room=87DY68
			go func() {
				err := g.netClient.Connect(serverAddress) // 172.20.10.2
				if err != nil {
					log.Println("Connection failed:", err)
				} else {
					g.state = sPlayMenu
				}
			}()
		}
		if g.menu.BtnQuit.IsClicked() {
			g.audioManager.Play("click_button")
			return ebiten.Termination
		}
	case sPlayMenu:
		for i, room := range g.playMenu.Rooms {
			if room.Btn.IsClicked() {
				g.state = sWaitingGame
				go func() {
					err := g.netClient.JoinGame(room.Id, isBotGame)
					g.playMenu.RoomIndex = i
					if err != nil {
						g.state = sMainMenu
						log.Println("Connection failed:", err)
					}
				}()
				break
			}
		}
		g.state = sWaitingGame
		// TODO
	case sWaitingGame:
		g.waitingMenu.RotationAngle += 0.08

		if g.waitingMenu.RotationAngle > math.Pi*2 {
			g.waitingMenu.RotationAngle -= math.Pi * 2
		}
		if g.audioManager.IsPlaying("main_menu_music") {
			g.audioManager.Stop("main_menu_music")

			err := g.audioManager.LoadMusic("waiting_opponent_music", "client/assets/waiting_opponent.mp3")
			if err != nil {
				log.Println("Could not load music:", err)
			}
			g.audioManager.Play("waiting_opponent_music")
		}
	case sGamePlaying:
		if !g.isMyTurn {
			return nil
		}

		cellX, cellY, ok := g.grid.OnClick()
		if !ok {
			return nil
		}

		err := g.netClient.SendPacket(common.Packet{
			Type: common.MsgClick,
			Data: func() json.RawMessage {
				payload, _ := json.Marshal(common.ClickPayload{
					X: cellX,
					Y: cellY,
				})
				return payload
			}(),
		})
		if err != nil {
			log.Println("Failed to send move:", err)
		}
	}
	return nil
}

// Draw the game screen.
// Called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case sMainMenu:
		ui.RenderMenu(screen, g.menu)
	case sWaitingGame:
		ui.RenderWaitingGame(screen, g.waitingMenu)
	case sPlayMenu:
		ui.RenderPlayMenu(screen, g.playMenu)
	case sGamePlaying:
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	case sGameWin:
		ui.RenderWin(screen)
	case sGameLose:
		ui.RenderLose(screen)
	case sGameDraw:
		ui.RenderDraw(screen)
	}
}

// Defines the game's screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ui.WindowWidth, ui.WindowHeight
}

// Handles incoming network messages.
func (g *Game) handleNetwork() {
	if g.netClient == nil {
		return
	}

	for {
		packet := g.netClient.Poll()
		if packet == nil {
			break // No more messages
		}

		switch packet.Type {
		case common.MsgRooms:
			var p common.RoomsPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			for roomId, playerCount := range p.Rooms {
				log.Printf("%d %s", playerCount, roomId)
				room := &ui.Room{Id: roomId, PlayerCount: playerCount}
				g.playMenu.Rooms = append(g.playMenu.Rooms, *room)
			}

		case common.MsgGameStart:
			var p common.GameStartPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			g.mySymbol = p.YouAre
			g.state = sGamePlaying // Server authorized us to start
			log.Printf("Game Started! I am Player %d", g.mySymbol)

		case common.MsgUpdate:
			var p common.UpdatePayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			g.grid.BoardData = p.Board
			g.isMyTurn = (p.Turn == g.mySymbol)
			log.Println("Board updated")

		case common.MsgGameOver:
			var p common.GameOverPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			if p.Winner == g.mySymbol {
				g.state = sGameWin
				log.Println("You Win!")
			} else if p.Winner == common.Empty {
				g.state = sGameDraw
				log.Println("It's a Draw!")
			} else {
				g.state = sGameLose
				log.Println("You Lose!")
			}
		}
	}
}

// Initialize UI elements like menus, grid, etc.
func (g *Game) initUIElements() {
	g.menu = ui.NewMainMenu()
	g.playMenu = &ui.PlayMenu{}
	g.waitingMenu = &ui.WaitingMenu{}
	g.grid = &ui.Grid{
		Col: ui.GridCol,
	}
}

func (g *Game) initAudio() {
	g.audioManager = audio.NewAudioManager()
	err := g.audioManager.LoadSound("click_button", "client/assets/click_button.wav")
	if err != nil {
		log.Printf("Error loading sound: %v", err)
	}
}
