package main

import (
	"Goonker/client/audio"
	"Goonker/client/ui"
	"Goonker/common"
	"fmt"
	"math/rand"
	"time"

	"encoding/json"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// States of the application
	sInit = iota
	sMainMenu
	sRoomsMenu
	sWaitingGame
	sGamePlaying
	sChallenge
	sGameWin
	sGameLose
	sGameDraw

	// Network configuration
	serverAddress = "wss://goonker.saikoon.ch/ws" // goonker.saikoon.ch
	isBotGame     = true
)

type Game struct {
	menu          *ui.MainMenu
	roomsMenu     *ui.RoomsMenu
	waitingMenu   *ui.WaitingMenu
	challengeMenu *ui.ChallengeMenu
	state         int
	netClient     *NetworkClient
	grid          *ui.Grid
	audioManager  *audio.AudioManager

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
	err := g.audioManager.LoadMusic("main_menu_music", "main_menu.mp3")
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
					g.state = sRoomsMenu
					err := g.netClient.GetRooms()
					if err != nil {
						log.Println("Could not get rooms : ", err)
					}
				}
			}()
		}
		if g.menu.BtnQuit.IsClicked() {
			g.audioManager.Play("click_button")
			return ebiten.Termination
		}
	case sRoomsMenu:
		g.roomsMenu.RoomField.Update()

		// Back to main menu
		if g.roomsMenu.BtnBack.IsClicked() {
			g.audioManager.Play("click_button")
			g.netClient.Disconnect()
			g.state = sMainMenu
		}

		// Join a bot game with a specialized ID
		if g.roomsMenu.BtnPlayBot.IsClicked() {
			g.audioManager.Play("click_button")
			// Create a bot game with a specialized ID
			err := g.netClient.JoinGame(fmt.Sprintf("BOT_%d", time.Now().Unix()/1000), true)
			if err != nil {
				log.Println("Connection failed:", err)
			}
			g.state = sWaitingGame
		}

		// Create a room with a random ID.
		if g.roomsMenu.BtnCreateRoom.IsClicked() {
			g.audioManager.Play("click_button")
			newRoomId := fmt.Sprintf("%d", time.Now().Unix()/1000)
			err := g.netClient.JoinGame(newRoomId, false)
			if err != nil {
				log.Println("Connection failed:", err)
			}
			g.state = sWaitingGame
			g.waitingMenu.RoomId = newRoomId
		}

		roomsNbr := len(g.roomsMenu.Rooms)
		if roomsNbr <= 0 {
			break
		}

		if g.roomsMenu.BtnJoinGame.IsClicked() {

			// Join the selected room
			roomId := g.roomsMenu.RoomField.Text

			// Join a random room if nothing was given
			if roomId == "" {
				roomIndex := rand.Intn(roomsNbr)
				roomId = g.roomsMenu.Rooms[roomIndex].Id
			}

			err := g.netClient.JoinGame(roomId, false)
			// g.roomsMenu.RoomIndex = roomIndex
			if err != nil {
				log.Println("Connection failed:", err)
			}
			g.state = sWaitingGame
			g.waitingMenu.RoomId = roomId
			break
		}

		// Join an existing room
		for i, room := range g.roomsMenu.Rooms {
			if room.JoinBtn.IsClicked() {
				err := g.netClient.JoinGame(room.Id, false)
				g.roomsMenu.RoomIndex = i
				if err != nil {
					log.Println("Connection failed:", err)
				}
				g.state = sWaitingGame
				break
			}
		}
	case sWaitingGame:
		g.waitingMenu.RotationAngle += 0.08

		if g.waitingMenu.RotationAngle > math.Pi*2 {
			g.waitingMenu.RotationAngle -= math.Pi * 2
		}
		if g.audioManager.IsPlaying("main_menu_music") {
			g.audioManager.Stop("main_menu_music")

			err := g.audioManager.LoadMusic("waiting_opponent_music", "waiting_opponent.mp3")
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
		g.audioManager.Play("place_symbol")
		err := g.netClient.PlaceSymbol(cellX, cellY)
		if err != nil {
			log.Println(err)
		}
	case sChallenge:
		g.challengeMenu.Clock.Update()
		for i, ansBtn := range g.challengeMenu.Answers {
			if ansBtn.IsClicked() {
				g.audioManager.Play("challenge")
				err := g.netClient.AnswerChallenge(i)
				if err != nil {
					log.Println("Connection failed:", err)
				}
				g.state = sGamePlaying
			}
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
	case sRoomsMenu:
		ui.RenderRoomsMenu(screen, g.roomsMenu)
	case sGamePlaying:
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	case sChallenge:
		ui.RenderChallenge(screen, g.challengeMenu)
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
			// Clear existing rooms
			g.roomsMenu.Rooms = nil

			// Update rooms list
			for _, roomId := range p.Rooms {
				log.Printf("%s", roomId)
				// Initialize
				room := ui.NewRoom(roomId)
				g.roomsMenu.Rooms = append(g.roomsMenu.Rooms, room)
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
			// Ensure the game state
			g.state = sGamePlaying

			var p common.UpdatePayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			g.grid.BoardData = p.Board
			g.isMyTurn = (p.Turn == g.mySymbol)
			log.Println("Board updated")

		case common.MsgChallenge:
			var payload common.ChallengePayload
			if err := json.Unmarshal(packet.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			g.challengeMenu = ui.NewChallengeMenu(payload)
			g.state = sChallenge

			g.challengeMenu.Clock = *ui.NewTimer(common.ChallengeTime * time.Second)
			g.challengeMenu.Clock.OnEnd = func() {
				g.audioManager.Play("challenge")
				g.state = sGamePlaying
				err := g.netClient.AnswerChallenge(-1)
				if err != nil {
					log.Println("Connection failed:", err)
				}
			}
		case common.MsgGameOver:
			var p common.GameOverPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			switch p.Winner {
			case g.mySymbol:
				g.state = sGameWin
				if g.audioManager.IsPlaying("challenge") {
					g.audioManager.Stop("challenge")
				}
				g.audioManager.Play("win")
				log.Println("You Win!")
			case common.Empty:
				g.state = sGameDraw
				if g.audioManager.IsPlaying("challenge") {
					g.audioManager.Stop("challenge")
				}
				g.audioManager.Play("lose")
				log.Println("It's a Draw!")
			default:
				g.state = sGameLose
				if g.audioManager.IsPlaying("challenge") {
					g.audioManager.Stop("challenge")
				}
				g.audioManager.Play("lose")
				log.Println("You Lose!")
			}
		}
	}
}

// Initialize UI elements like menus, grid, etc.
func (g *Game) initUIElements() {
	g.menu = ui.NewMainMenu()
	g.roomsMenu = ui.NewRoomsMenu()
	g.waitingMenu = &ui.WaitingMenu{}
	g.grid = &ui.Grid{
		Col: ui.GridCol,
	}
}

// Initialize audio manager and load sounds
func (g *Game) initAudio() {
	g.audioManager = audio.NewAudioManager()
	err := g.audioManager.LoadSound("click_button", "click_button.wav")
	if err != nil {
		log.Printf("Error loading sound: %v", err)
	}
	err = g.audioManager.LoadSound("place_symbol", "place_symbol.wav")
	if err != nil {
		log.Printf("Error loading music: %v", err)
	}
	err = g.audioManager.LoadSound("win", "win.wav")
	if err != nil {
		log.Printf("Error loading music: %v", err)
	}
	err = g.audioManager.LoadSound("lose", "lose.wav")
	if err != nil {
		log.Printf("Error loading music: %v", err)
	}
	err = g.audioManager.LoadSound("challenge", "challenge.wav")
	if err != nil {
		log.Printf("Error loading music: %v", err)
	}
}
