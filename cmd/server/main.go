package main

import (
	"log"
	"net/http"

	"github.com/Bitlatte/p2p-chat-server/internal/server"
)

func main() {
	s := server.NewServer()
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}