package hub

import (
	"sync"
)

type Hub struct {
	rooms map[string]*Room
	mutex sync.Mutex
}

// Singleton Global Hub
var GlobalHub = &Hub{
	rooms: make(map[string]*Room),
}

func (h *Hub) GetRoom(roomID string) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.rooms[roomID]
}

func (h *Hub) CreateRoom(roomID string, isBot bool) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		return room
	}

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

// GetAvailableRooms returns a map of the rooms IDs and the amount of players in it
// If the room is full, it is not included in the result
func (h *Hub) GetAvailableRooms() map[string]int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	availableRooms := make(map[string]int)
	for roomID := range h.rooms {
		if !h.rooms[roomID].IsFull() {
			playerCount := len(h.rooms[roomID].Players)
			availableRooms[roomID] = playerCount
		}
	}
	return availableRooms
}
