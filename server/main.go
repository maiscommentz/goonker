package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Goonker/common"
	"Goonker/server/hub"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Only for local testing
	})
	if err != nil {
		log.Printf("Error upgrading websocket: %v", err)
		return
	}

	// Handshake: Wait for the "join" message
	// We give ourselves a timeout of 5 seconds to receive the join, otherwise we cut off.
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var packet common.Packet
	if err := wsjson.Read(ctx, c, &packet); err != nil {
		log.Printf("Connection closed before join: %v", err)
		c.Close(websocket.StatusPolicyViolation, "Expected Join Packet")
		return
	}

	// Validation of message type
	if packet.Type != common.MsgJoin {
		log.Printf("First packet was not join: %s", packet.Type)
		c.Close(websocket.StatusPolicyViolation, "First message must be 'join'")
		return
	}

	// Parsing of Payload
	var joinData common.JoinPayload
	if err := json.Unmarshal(packet.Data, &joinData); err != nil {
		log.Printf("Invalid join payload: %v", err)
		c.Close(websocket.StatusProtocolError, "Invalid Payload")
		return
	}

	if joinData.RoomID == "" {
		c.Close(websocket.StatusPolicyViolation, "Room ID required")
		return
	}

	// Creation / Retrieval of the Room via the Hub
	// Note: The Hub manages synchronization, no need for a mutex here.
	room := hub.GlobalHub.CreateRoom(joinData.RoomID, joinData.IsBot)

	log.Printf("Client joining room '%s' (Bot: %v)", joinData.RoomID, joinData.IsBot)

	// Add the player to the Room
	// From this point, the Room manages the connection (read loop)
	// We cancel the handshake timeout context because the connection is now long-lived
	pid := room.AddPlayer(c)

	if pid == common.Empty {
		log.Println("Room is full, rejecting client")
		// On ferme proprement si la salle est pleine
		c.Close(websocket.StatusPolicyViolation, "Room is full")
	} else {
		log.Printf("Player assigned ID: %d in room %s", pid, joinData.RoomID)
	}
}