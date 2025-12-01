package hub

import (
	"context"
	"log"
	"sync"

	"nhooyr.io/websocket"
)

type Hub struct {
	roomsMu sync.RWMutex
	rooms   map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

// HandleJoin finds the room and adds the player
func (h *Hub) HandleJoin(roomID string, conn *websocket.Conn, ctx context.Context) {
	h.roomsMu.Lock()
	room, exists := h.rooms[roomID]
	if !exists {
		// Create new room if it doesn't exist
		room = NewRoom(roomID)
		h.rooms[roomID] = room
		log.Printf("Created new room: %s", roomID)
	}
	h.roomsMu.Unlock()

	// Add player to the room (Blocks until client disconnects)
	room.Join(conn, ctx)
	
	h.checkAndDestroy(roomID, room)
}

// checkAndDestroy removes the room from the map if it is empty
func (h *Hub) checkAndDestroy(roomID string, room *Room) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	// Critical: Check if the room in the map is actually the same instance
	// (Prevents race conditions if a room was deleted and recreated instantly)
	if existingRoom, ok := h.rooms[roomID]; ok && existingRoom == room {
		if !room.HasPlayers() {
			delete(h.rooms, roomID)
			log.Printf("[%s] Room empty and deleted", roomID)
		} else {
			log.Printf("[%s] Player left, but room still active", roomID)
		}
	}
}