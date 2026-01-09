package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Goonker/common"
	"Goonker/server/hub"
	"Goonker/server/utils"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Server configuration constants
const (
	// Network configuration
	ServerPort       = ":8080"
	WsRoute          = "/ws"
	HandshakeTimeout = 5 * time.Second

	// Closure Reasons
	ErrExpectedJoin    = "Expected Join Packet"
	ErrFirstMustBeJoin = "First message must be 'join'"
	ErrInvalidPayload  = "Invalid Payload"
	ErrRoomIDRequired  = "Room ID required"
	ErrRoomFull        = "Room is full"
)

// main is the entry point of the server application.
func main() {
	// Register the WebSocket handler
	http.HandleFunc(WsRoute, wsHandler)

	// Load the challenges
	utils.LoadChallenges()

	// Start the server
	log.Printf("Starting server on port %s...", ServerPort)
	if err := http.ListenAndServe(ServerPort, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// wsHandler handles the initial HTTP upgrade and the application-layer handshake.
// Once the player is validated, control is passed to the Hub/Room.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket (skip verify for local dev)
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("Error upgrading websocket: %v", err)
		return
	}

	// Context for the connection lifecycle while in the lobby/handshake phase
	// We use the request context which is cancelled when the connection closes
	ctx := r.Context()

	for {
		// Read a packet
		var packet common.Packet
		if err := wsjson.Read(ctx, c, &packet); err != nil {
			log.Printf("Connection closed or error reading packet: %v", err)
			return
		}

		switch packet.Type {
		case common.MsgJoin:
			// Parse the Join payload
			var joinData common.JoinPayload
			if err := json.Unmarshal(packet.Data, &joinData); err != nil {
				log.Printf("Invalid join payload: %v", err)
				err = c.Close(websocket.StatusProtocolError, ErrInvalidPayload)
				if err != nil {
					log.Println(err)
				}

				return
			}

			// Validate RoomID presence
			if joinData.RoomID == "" {
				err = c.Close(websocket.StatusPolicyViolation, ErrRoomIDRequired)
				if err != nil {
					log.Println(err)
				}
				return
			}

			// Let the Hub assign the player to a new or existing room
			room := hub.GlobalHub.CreateRoom(joinData.RoomID, joinData.IsBot)
			log.Printf("Client joining room '%s' (Bot: %v)", joinData.RoomID, joinData.IsBot)
			pid := room.AddPlayer(c)

			// Validation of assigned PlayerID, otherwise room is full
			if pid == common.Empty {
				log.Println("Room is full, rejecting client")
				err = c.Close(websocket.StatusPolicyViolation, ErrRoomFull)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Printf("Player assigned ID: %d in room %s", pid, joinData.RoomID)
			}

			// Once joined, the Room takes over the connection (reading/writing)
			// so we must exit this handler loop to avoid concurrent reading.
			return

		case common.MsgGetRooms:
			// Fetch available rooms from the Hub
			rooms := hub.GlobalHub.GetAvailableRooms()

			// Send the list back to the client
			payload := common.RoomsPayload{Rooms: rooms}
			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshaling rooms payload: %v", err)
				continue
			}

			response := common.Packet{
				Type: common.MsgRooms,
				Data: data,
			}

			// Use a timeout for writing
			writeCtx, cancel := context.WithTimeout(ctx, HandshakeTimeout)
			defer cancel()

			if err := wsjson.Write(writeCtx, c, response); err != nil {
				log.Printf("Error sending MsgRooms: %v", err)
				return
			}

		default:
			// State machine is permissive in 'lobby', ignore unexpected messages
			log.Printf("Unexpected message type: %s", packet.Type)
		}
	}
}
