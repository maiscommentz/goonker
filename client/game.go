package main

import (
	"Goonker/client/ui"
	"Goonker/common"
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// States of the application
	sInit = iota
	sMenu
	sPlayMenu
	sGamePlaying
	sGameWin
	sGameLose
)

type Game struct {
	menu      *ui.Menu
	playMenu  *PlayMenu
	state     int
	netClient *NetworkClient
	grid      *ui.Grid

	mySymbol common.PlayerID // 1 for X, 2 for O
	isMyTurn bool
}

type PlayMenu struct {
	// TODO: Add fields for play menu (e.g., room selection, bot option)
}

/**
 * Initializes the game state.
 */
func (g *Game) Init() {
	// Initialize network client
	g.netClient = NewNetworkClient()

	// Initialize menu
	g.menu = &ui.Menu{}

	// Center buttons
	buttonWidth, buttonHeight := 200.0, 60.0
	centerX := (float64(ui.WindowWidth) - buttonWidth) / 2

	// Create buttons
	g.menu.BtnPlay = ui.NewButton(centerX, 200, buttonWidth, buttonHeight, "Play")
	g.menu.BtnQuit = ui.NewButton(centerX, 300, buttonWidth, buttonHeight, "Quit")

	// Pre-render menu image
	if g.menu.MenuImage == nil {
		img := ui.DrawMenu(ui.WindowWidth, ui.WindowHeight, ui.GameTitle)
		g.menu.MenuImage = ebiten.NewImageFromImage(img)
	}

	// Initialize the grid
	g.grid = ui.NewGrid(ui.GridCol, ui.GridCol)

	// Set initial state
	g.state = sMenu
}

/**
 * Updates the game state every tick.
 * Typically called every tick (1/60[s] by default).
 */
func (g *Game) Update() error {

	// Always poll the network for incoming messages first
	g.handleNetwork()

	switch g.state {
	case sInit:
		g.Init()
	case sMenu:
		if g.menu.BtnPlay.IsClicked() {
			// TODO: This block will be placed in PlayMenu later
			// Try to connect to server (Async)
			// Note: For WASM/Localhost testing use ws://localhost:8080/ws?room=87DY68
			go func() {
				err := g.netClient.Connect("ws://localhost:8080/ws", "87DY68", false) // 172.20.10.2
				if err != nil {
					log.Println("Connection failed:", err)
				}
			}()
			// TODO: g.state = sPlayMenu
		}
		if g.menu.BtnQuit.IsClicked() {
			return ebiten.Termination
		}
	case sPlayMenu:
		//TODO: Handle Play Menu interactions
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

/**
 * Draws the game screen.
 * Called every frame (typically 1/60[s] for 60Hz display).
 */
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case sMenu:
		ui.RenderMenu(screen, g.menu)
	case sPlayMenu:
		//TODO: ui.RenderPlayMenu(...)
	case sGamePlaying:
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	case sGameWin:
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	case sGameLose:
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	}
}

/**
 * Defines the game's screen dimensions.
 */
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ui.WindowWidth, ui.WindowHeight
}

/**
 * Handles incoming network messages.
 */
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
			json.Unmarshal(packet.Data, &p)

			g.mySymbol = p.YouAre
			g.state = sGamePlaying // Server authorized us to start
			log.Printf("Game Started! I am Player %d", g.mySymbol)

		case common.MsgUpdate:
			var p common.UpdatePayload
			json.Unmarshal(packet.Data, &p)

			g.grid.BoardData = p.Board
			g.isMyTurn = (p.Turn == g.mySymbol)
			log.Println("Board updated")
		}
	}
}
