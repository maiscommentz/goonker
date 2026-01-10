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

// Constants
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
	serverAddress = "wss://goonker.saikoon.ch/ws"
	isBotGame     = true
)

// Game represents the game state
type Game struct {
	menu          *ui.MainMenu
	roomsMenu     *ui.RoomsMenu
	waitingMenu   *ui.WaitingMenu
	challengeMenu *ui.ChallengeMenu
	gameOverMenu  *ui.GameOverMenu
	state         int
	netClient     *NetworkClient
	grid          *ui.Grid
	audioManager  *audio.AudioManager

	mySymbol common.PlayerID // 1 for X, 2 for O
	isMyTurn bool
}

// Init the game
func (g *Game) Init() {
	// Initialize network client which handles communication with the server
	g.netClient = NewNetworkClient()

	// Initialize the UI elements and assets
	ui.Init()

	// Initialize UI elements specific to game states (menus, grid)
	g.initUIElements()

	// Initialize Audio Manager and load game sounds
	g.initAudio()

	// Play main menu music
	err := g.audioManager.LoadMusic("main_menu_music", "main_menu.mp3")
	if err != nil {
		log.Println("Could not load music:", err)
	}
	g.audioManager.Play("main_menu_music")

	// Set initial state to Main Menu
	g.state = sMainMenu
}

// Update the game state every tick.
// Typically called every tick (1/60[s] by default).
func (g *Game) Update() error {

	// Always poll the network for incoming messages first
	g.handleNetwork()

	switch g.state {
	case sInit:
		// Initialize the game if in Init state
		g.Init()
	case sMainMenu:
		// Handle Main Menu interactions

		// Click on play
		if g.menu.BtnPlay.IsClicked() {
			g.audioManager.Play("click_button")
			// Try to connect to server (Async)
			go func() {
				err := g.netClient.Connect(serverAddress)
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
		// Click on quit
		if g.menu.BtnQuit.IsClicked() {
			g.audioManager.Play("click_button")
			return ebiten.Termination
		}
	case sRoomsMenu:
		// Handle Rooms Menu interactions

		// Update the room field for text input
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
			// Create a bot game with a specialized ID based on timestamp
			err := g.netClient.JoinGame(fmt.Sprintf("BOT_%d", time.Now().Unix()), true)
			if err != nil {
				log.Println("Connection failed:", err)
			}
			g.state = sWaitingGame
		}

		// Create a room with a random ID.
		if g.roomsMenu.BtnCreateRoom.IsClicked() {
			g.audioManager.Play("click_button")
			newRoomId := fmt.Sprintf("%d", time.Now().Unix())
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

		// Click on join game
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

		// Join an existing room from the list
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
		// Handle Waiting Screen animations and logic

		// Update the animation wheel
		g.waitingMenu.RotationAngle += 0.08

		if g.waitingMenu.RotationAngle > math.Pi*2 {
			g.waitingMenu.RotationAngle -= math.Pi * 2
		}

		// Play the waiting music if not already playing
		if g.audioManager.IsPlaying("main_menu_music") {
			g.audioManager.Stop("main_menu_music")

			err := g.audioManager.LoadMusic("waiting_opponent_music", "waiting_opponent.mp3")
			if err != nil {
				log.Println("Could not load music:", err)
			}
			g.audioManager.Play("waiting_opponent_music")
		}
	case sGamePlaying:
		// Handle Game Playing state

		// Click on a cell
		if !g.isMyTurn {
			return nil
		}

		// Check for grid clicks
		cellX, cellY, ok := g.grid.OnClick()
		if !ok {
			return nil
		}
		g.audioManager.Play("place_symbol")
		// Send move to server
		err := g.netClient.PlaceSymbol(cellX, cellY)
		if err != nil {
			log.Println(err)
		}
	case sChallenge:
		// Handle Challenge state (Mini-game/Quiz)

		// Update the clock
		g.challengeMenu.Clock.Update()
		for i, ansBtn := range g.challengeMenu.Answers {
			if ansBtn.IsClicked() {
				g.audioManager.Play("challenge")
				// Send answer to server
				err := g.netClient.AnswerChallenge(i)
				if err != nil {
					log.Println("Connection failed:", err)
				}
				g.state = sGamePlaying
			}
		}
	case sGameWin, sGameLose, sGameDraw:
		// Handle Game Over states

		// Click on back
		if g.gameOverMenu.BtnBack.IsClicked() {
			g.audioManager.Play("click_button")
			g.netClient.Disconnect()
			// Reconnect to lobby
			go func() {
				err := g.netClient.Connect(serverAddress)
				if err != nil {
					log.Println("Connection failed:", err)
					g.state = sMainMenu
				} else {
					g.state = sRoomsMenu
					err := g.netClient.GetRooms()
					if err != nil {
						log.Println("Could not get rooms : ", err)
					}
				}
			}()
		}
	}
	return nil
}

// Draw the game screen.
// Called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case sMainMenu:
		// Draw Main Menu
		ui.RenderMenu(screen, g.menu)
	case sWaitingGame:
		// Draw Waiting Screen
		ui.RenderWaitingGame(screen, g.waitingMenu)
	case sRoomsMenu:
		// Draw Rooms List
		ui.RenderRoomsMenu(screen, g.roomsMenu)
	case sGamePlaying:
		// Draw Game Board
		ui.RenderGame(screen, g.grid, g.isMyTurn)
	case sChallenge:
		// Draw Challenge/Quiz Interface
		ui.RenderChallenge(screen, g.challengeMenu)
	case sGameWin:
		// Draw Win Screen
		ui.RenderWin(screen, g.gameOverMenu)
	case sGameLose:
		// Draw Lose Screen
		ui.RenderLose(screen, g.gameOverMenu)
	case sGameDraw:
		// Draw Draw Screen
		ui.RenderDraw(screen, g.gameOverMenu)
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
		// Non-blocking poll for new packets
		packet := g.netClient.Poll()
		if packet == nil {
			break // No more messages
		}

		switch packet.Type {
		case common.MsgRooms:
			// Handle room list update
			var p common.RoomsPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}
			// Clear existing rooms
			g.roomsMenu.Rooms = nil

			// Update rooms list with new data
			for _, roomId := range p.Rooms {
				log.Printf("%s", roomId)
				// Initialize new Room UI element
				room := ui.NewRoom(roomId)
				g.roomsMenu.Rooms = append(g.roomsMenu.Rooms, room)
			}

		case common.MsgGameStart:
			// Handle game start signal
			var p common.GameStartPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			g.mySymbol = p.YouAre
			g.state = sGamePlaying // Server authorized us to start
			log.Printf("Game Started! I am Player %d", g.mySymbol)

		case common.MsgUpdate:
			// Handle board update
			// Ensure the game state is set to playing (recovers from network lag/missed packets)
			g.state = sGamePlaying

			var p common.UpdatePayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			// Update local grid data
			g.grid.BoardData = p.Board
			g.isMyTurn = (p.Turn == g.mySymbol)
			log.Println("Board updated")

		case common.MsgChallenge:
			// Handle challenge trigger (quiz)
			var payload common.ChallengePayload
			if err := json.Unmarshal(packet.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			// Initialize and display challenge menu
			g.challengeMenu = ui.NewChallengeMenu(payload)
			g.state = sChallenge

			// Start challenge timer
			g.challengeMenu.Clock = *ui.NewTimer(common.ChallengeTime * time.Second)
			g.challengeMenu.Clock.OnEnd = func() {
				// Handle timer expiration
				g.audioManager.Play("challenge")
				g.state = sGamePlaying
				// Send empty answer (-1) on timeout
				err := g.netClient.AnswerChallenge(-1)
				if err != nil {
					log.Println("Connection failed:", err)
				}
			}
		case common.MsgGameOver:
			// Handle game over result
			var p common.GameOverPayload
			if err := json.Unmarshal(packet.Data, &p); err != nil {
				log.Printf("Failed to unmarshal %s: %v", packet.Type, err)
				continue
			}

			// Determine result and switch state/music
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
	// Initialize Main Menu
	g.menu = ui.NewMainMenu()
	// Initialize Rooms Menu
	g.roomsMenu = ui.NewRoomsMenu()
	// Initialize Game Over Menu
	g.gameOverMenu = ui.NewGameOverMenu()
	// Initialize Waiting Menu
	g.waitingMenu = &ui.WaitingMenu{}
	// Initialize Game Grid with default columns
	g.grid = &ui.Grid{
		Col: ui.GridCol,
	}
}

// Initialize audio manager and load sounds
func (g *Game) initAudio() {
	g.audioManager = audio.NewAudioManager()

	// Load sound effects
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
