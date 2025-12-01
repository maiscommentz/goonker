package main

import (
	"encoding/json"
	"Goonker/client/ui"
	"Goonker/common"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	boardImage *ebiten.Image
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
	menu      *Menu
	playMenu  *PlayMenu
	state     int
	netClient *NetworkClient

	mySymbol  common.PlayerID // 1 for X, 2 for O
	boardData [3][3]common.PlayerID
	isMyTurn   bool
}

type Menu struct {
	menuImage *ebiten.Image
	btnPlay   *ui.Button
	btnQuit   *ui.Button
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
	g.menu = &Menu{}

	// Center buttons
	buttonWidth, buttonHeight := 200.0, 60.0
	centerX := (float64(WindowWidth) - buttonWidth) / 2

	// Create buttons
	g.menu.btnPlay = ui.NewButton(centerX, 200, buttonWidth, buttonHeight, "Play")
	g.menu.btnQuit = ui.NewButton(centerX, 300, buttonWidth, buttonHeight, "Quit")

	// Pre-render menu image
	if g.menu.menuImage == nil {
		img := ui.DrawMenu(WindowWidth, WindowHeight, GameTitle)
		g.menu.menuImage = ebiten.NewImageFromImage(img)
	}

	// Pre-render board image
	if boardImage == nil {
		grid := ui.DrawGrid()
		boardImage = ebiten.NewImageFromImage(grid)
	}

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
		if g.menu.btnPlay.IsClicked() {
			// TODO: This block will be placed in PlayMenu later
			// Try to connect to server (Async)
			// Note: For WASM/Localhost testing use ws://localhost:8080/ws?room=87DY68
			go func() {
				err := g.netClient.Connect("ws://localhost:8080/ws", "87DY68", false)
				if err != nil {
					log.Println("Connection failed:", err)
				}
			}()
			// TODO: g.state = sPlayMenu
		}
		if g.menu.btnQuit.IsClicked() {
			return ebiten.Termination
		}
	case sPlayMenu:
		//TODO: Handle Play Menu interactions
	case sGamePlaying:
		if g.isMyTurn && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Still TODO: Handle click on board
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
		ui.RenderMenu(screen, g.menu.menuImage, g.menu.btnPlay, g.menu.btnQuit)
	case sPlayMenu:
		//TODO: ui.RenderPlayMenu(...)
	case sGamePlaying:
		ui.RenderGame(screen, boardImage)
	case sGameWin:
		ui.RenderGame(screen, boardImage)
	case sGameLose:
		ui.RenderGame(screen, boardImage)
	}
}

/**
 * Defines the game's screen dimensions.
 */
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
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
			
			g.boardData = p.Board
			g.isMyTurn = (p.Turn == g.mySymbol)
			log.Println("Board updated")
		}
	}
}