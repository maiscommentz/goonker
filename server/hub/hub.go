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

// CleanUp supprime les rooms vides (à appeler périodiquement si besoin)
func (h *Hub) RemoveRoom(roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.rooms, roomID)
}