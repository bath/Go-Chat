package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message a structure for the browser and Go to use json objects
type Message struct {
	Email	string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

// dict of connected clients .. 
// the client map's key is a pointer to a websocket connection and the value is a bool
var clients = make(map[*websocket.Conn]bool) 
var broadcast = make(chan Message) // broadcast channel of type Message

var upgrader = websocket.Upgrader{}

func main() {
	// Create simple file server
	log.Print("Creating file server handler")
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket connection establishments
	log.Print("Starting websocket connection handler")
	http.HandleFunc("/ws", handleConnections)

	// Listen for inbound chat messages
	go handleMessages()

	// start server on localhost port 8000, log errors if present
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Convert initial GET request into a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("GETtoWebsocket: ", err)
	}

	defer ws.Close()

	// log in the client map that this websocket is currently active
	clients[ws] = true

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		
		// write the message to the broadcast channel to output on all ws conns
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <- broadcast
		
		for client := range clients { // grab the key, value(undeclared == ignored) for each client in clients
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}