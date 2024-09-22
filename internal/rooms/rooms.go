package rooms

import (
	"sync"
)

type messagePayload struct {
	SenderID string
	Message  []byte
}

type Room struct {
	Name         string
	Type         string             // "text" or "multimedia"
	Whitelist    map[string]bool    // UserID whitelist
	Participants map[string]*Client // Connected clients
	register     chan *Client
	unregister   chan *Client
	broadcast    chan messagePayload
	mu           sync.RWMutex
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) CreateRoom(name string, roomType string, whitelist []string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	wl := make(map[string]bool)
	for _, userID := range whitelist {
		wl[userID] = true
	}
	room := &Room{
		Name:         name,
		Type:         roomType,
		Whitelist:    wl,
		Participants: make(map[string]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan messagePayload),
	}
	rm.rooms[name] = room
	go room.Run()
}

func (r *Room) RegisterClient(c *Client) {
	r.register <- c
}

func (rm *RoomManager) GetRoom(name string) *Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[name]
}

// Check if a user is allowed to join a room
func (rm *RoomManager) CanJoinRoom(roomName, userID string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, exists := rm.rooms[roomName]
	if !exists {
		return false
	}
	return room.CanUserJoin(userID)
}

func (r *Room) CanUserJoin(userID string) bool {
	// Check if user is blacklisted (if you have implemented blacklisting)
	if r.IsUserBlacklisted(userID) {
		return false
	}
	// If whitelist is empty, room is public
	if len(r.Whitelist) == 0 {
		return true
	}
	// Check if user is whitelisted
	return r.Whitelist[userID]
}

// Additional methods for managing blacklists and whitelists
func (r *Room) IsUserBlacklisted(userID string) bool {
	// Implement blacklist logic if needed
	return false
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.mu.Lock()
			r.Participants[client.UserID] = client
			r.mu.Unlock()
		case client := <-r.unregister:
			r.mu.Lock()
			if _, ok := r.Participants[client.UserID]; ok {
				delete(r.Participants, client.UserID)
				close(client.Send)
			}
			r.mu.Unlock()
		case message := <-r.broadcast:
			r.mu.RLock()
			for userID, client := range r.Participants {
				if userID != message.SenderID {
					select {
					case client.Send <- message.Message:
					default:
						close(client.Send)
						delete(r.Participants, client.UserID)
					}
				}
			}
			r.mu.RUnlock()
		}
	}
}
