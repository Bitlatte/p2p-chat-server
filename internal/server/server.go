package server

import (
	"github.com/Bitlatte/p2p-chat-server/internal/rooms"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Rooms  *rooms.RoomManager
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		Rooms:  rooms.NewRoomManager(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// WebSocket endpoint for clients
	s.Router.HandleFunc("/ws", s.handleWebSocket)

	// Admin endpoints
	s.Router.HandleFunc("/admin/create-room", s.handleCreateRoom).Methods("POST")
	// s.Router.HandleFunc("/admin/blacklist-user", s.handleBlacklistUser).Methods("POST")
	// s.Router.HandleFunc("/admin/whitelist-room", s.handleWhitelistRoom).Methods("POST")
}
