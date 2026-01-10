package hub

import (
	"sync"
)

// Hub represents the hub that manages rooms
type Hub struct {
	rooms map[string]*Room
	mutex sync.Mutex
}

// Singleton Global Hub
var GlobalHub = &Hub{
	rooms: make(map[string]*Room),
}

// GetRoom returns a room by its ID
func (h *Hub) GetRoom(roomID string) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.rooms[roomID]
}

// CreateRoom creates a new room if it doesn't exist
func (h *Hub) CreateRoom(roomID string, isBot bool) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// If the room already exists, return it
	if room, exists := h.rooms[roomID]; exists {
		return room
	}

	// Create the room
	newRoom := NewRoom(roomID, isBot)
	h.rooms[roomID] = newRoom
	return newRoom
}

// RemoveRoom deletes a room from the hub
func (h *Hub) RemoveRoom(roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.rooms, roomID)
}

// GetAvailableRooms returns a slice of the room IDs
// If the room is full, it is not included in the result
func (h *Hub) GetAvailableRooms() []string {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var availableRooms []string
	for roomID := range h.rooms {
		if !h.rooms[roomID].IsFull() {
			availableRooms = append(availableRooms, roomID)
		}
	}
	return availableRooms
}
