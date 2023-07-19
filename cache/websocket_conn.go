package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
)

// Function to broadcast the file changes to all connected clients via WebSocket.
func broadcastFileChanges(broadcastChannel <-chan []*FileChange, connectedClients map[*websocket.Conn]bool) {
	for fileChanges := range broadcastChannel {
		// Implement code to synchronize changes with other participants based on the received message.
		fmt.Println("Received message on channel 'cache_channel':", fileChanges)

		// Broadcast the file changes to all connected clients via WebSocket.
		for client := range connectedClients {
			for _, fileChange := range fileChanges {
				// Assuming you have a helper function to marshal FileChange to JSON.
				jsonData, err := json.Marshal(fileChange)
				if err != nil {
					log.Println("Error converting FileChange to JSON:", err)
					continue
				}
				client.WriteMessage(websocket.TextMessage, jsonData)
			}
		}
	}
}

// Main function to handle WebSocket connections.
func handleWebSocketConnections() {
	// WebSocket Upgrader configuration
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true //TODO: You should perform proper origin checks in a production environment.
		},
	}

	// Connected WebSocket clients (you may want to store them in a more scalable way).
	connectedClients := make(map[*websocket.Conn]bool)

	// Channel to receive broadcast messages from the "cache_channel".
	broadcastChannel := make(chan []*FileChange)
	var wg sync.WaitGroup
	wg.Add(1)

	// Start a goroutine to handle incoming broadcast messages and send them to connected clients via WebSocket.
	go broadcastFileChanges(broadcastChannel, connectedClients)

	// HTTP handler for WebSocket connections.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection to WebSocket:", err)
			return
		}
		fmt.Println(conn)

		// Register the WebSocket client in the connectedClients map.
		connectedClients[conn] = true

		// Close the WebSocket connection and remove it from the connectedClients map when the client disconnects.
		defer func() {
			conn.Close()
			delete(connectedClients, conn)
		}()

		// You can implement further logic here, e.g., handling client messages sent via WebSocket.
		// For example, if a client sends a chat message, you can process it and send it to other clients.
	})

	// HTTP server to handle WebSocket connections.
	http.ListenAndServe(":8080", nil)

	// Wait for termination signal to gracefully shut down the service.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	<-terminate
	wg.Done()
}
