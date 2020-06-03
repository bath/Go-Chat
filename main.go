package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // dict of connected clients
var broadcast = make(chan Message) // broadcast channel

var upgrader = websocket.Upgrader{}

// Message - a structure for the browser and Go to use json objects
type Message struct {
	Email	string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func main() {
	// Create simple file server
	log.Print("Starting File Server")
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket connection establishments
	log.Print("Starting websocket connection handler")
	http.HandleFunc("/ws", handleConnections)

	// Listen for inbound chat messages
	go handleMessages()
}