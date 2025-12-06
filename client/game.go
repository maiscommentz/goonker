package main

import (
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
	//sPlayMenu
	sWaitingGame
	sGamePlaying
	sGameWin
	sGameLose
	sGameDraw

	// Network configuration
	serverAddress = "ws://goonker.saikoon.ch/ws" // goonker.saikoon.ch
	roomId        = "87DY68"
	isBotGame     = true
)

type Game struct {
	menu        *ui.MainMenu
	waitingMenu *ui.WaitingMenu
	playMenu    *ui.PlayMenu
	state       int
	netClient   *NetworkClient
	grid        *ui.Grid

	mySymbol common.PlayerID // 1 for X, 2 for O
	isMyTurn bool
}

// Init the game
func (g *Game) Init() {
	// Initialize network client
	g.netClient = NewNetworkClient()

	// Initialize the UI
	ui.Init()

	// Initialize the main menu
	g.menu = ui.NewMainMenu()

	g.waitingMenu = &ui.WaitingMenu{}

	// Initialize the grid
	g.grid = &ui.Grid{
		Col: ui.GridCol,
	}

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
			// TODO: This block will be placed in PlayMenu later
			// Try to connect to server (Async)
			// Note: For WASM/Localhost testing use ws://localhost:8080/ws?room=87DY68
			g.state = sWaitingGame
			go func() {
				err := g.netClient.Connect(serverAddress, roomId, isBotGame) // 172.20.10.2
				if err != nil {
					g.state = sMainMenu
					log.Println("Connection failed:", err)
				}
			}()
			// TODO: g.state = sPlayMenu
		}
		if g.menu.BtnQuit.IsClicked() {
			return ebiten.Termination
		}
	//case sPlayMenu:
	//TODO: Handle Play Menu interactions
	case sWaitingGame:
		g.waitingMenu.RotationAngle += 0.08

		if g.waitingMenu.RotationAngle > math.Pi*2 {
			g.waitingMenu.RotationAngle -= math.Pi * 2
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
	//case sPlayMenu:
	//TODO: ui.RenderPlayMenu(...)
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
