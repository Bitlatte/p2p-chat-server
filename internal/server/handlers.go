package server

import (
	"encoding/json"
	"net/http"

	"github.com/Bitlatte/p2p-chat-server/internal/admin"
	"github.com/Bitlatte/p2p-chat-server/internal/auth"
	"github.com/Bitlatte/p2p-chat-server/internal/rooms"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking
		return true
	},
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}

	// Read initial message with unique values
	_, message, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	var initData struct {
		Value1   string `json:"value1"`
		Value2   string `json:"value2"`
		RoomName string `json:"roomName"`
	}
	if err := json.Unmarshal(message, &initData); err != nil {
		conn.Close()
		return
	}

	userID := auth.GenerateUserID(initData.Value1, initData.Value2)

	// Check if user can join the room
	if !s.Rooms.CanJoinRoom(initData.RoomName, userID) {
		conn.WriteMessage(websocket.TextMessage, []byte("Access denied to the room"))
		conn.Close()
		return
	}

	// Get the room
	room := s.Rooms.GetRoom(initData.RoomName)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found"))
		conn.Close()
		return
	}

	// Create a new client
	client := &rooms.Client{
		UserID: userID,
		Conn:   conn,
		Room:   room,
		Send:   make(chan []byte),
	}

	// Register the client with the room
	room.RegisterClient(client)

	// Start reading and writing
	go client.WritePump()
	client.ReadPump()
}

// Authentication middleware for admin endpoints
func (s *Server) adminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !admin.IsAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Whitelist []string `json:"whitelist"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.Rooms.CreateRoom(req.Name, req.Type, req.Whitelist)
	w.WriteHeader(http.StatusCreated)
}
