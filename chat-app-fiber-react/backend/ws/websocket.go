// websocket.go

package ws

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex
var typingUsers = make(map[*websocket.Conn]time.Time)

const typingTimeout = 3 * time.Second

// HandleWebSocket connects a new WebSocket client
func HandleWebSocket(c *websocket.Conn) {
	// Log when a new client connects
	log.Println("New client connected")

	// Lock and add client to the clients map
	mu.Lock()
	clients[c] = true
	mu.Unlock()

	typingTicker := time.NewTicker(1 * time.Second)
	defer typingTicker.Stop()

	go func() {
		for range typingTicker.C {
			mu.Lock()
			// If user is not typing anymore (timeout expired), remove them from the typing list
			if time.Since(typingUsers[c]) > typingTimeout {
				// Notify other clients that this user stopped typing
				for client := range clients {
					if err := client.WriteMessage(websocket.TextMessage, []byte("stop-typing:"+c.RemoteAddr().String())); err != nil {
						log.Println("Error sending stop typing message:", err)
					}
				}
				delete(typingUsers, c)
			}
			mu.Unlock()
		}
	}()

	// Read and broadcast messages
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Check if the message is a typing event
		if string(msg) == "typing" {
			// Mark this usre as typing
			mu.Lock()
			typingUsers[c] = time.Now()
			mu.Unlock()

			// Notify other clients that this user is typing
			for client := range clients {
				if err := client.WriteMessage(websocket.TextMessage, []byte("typing:"+c.RemoteAddr().String())); err != nil {
					log.Println("Error sending typing message:", err)
				}
			}
		}

		// Broadcast message to all other clients
		mu.Lock()
		for client := range clients {
			if err := client.WriteMessage(msgType, msg); err != nil {
				log.Println("Error sending message:", err)
			}
		}
		mu.Unlock()
	}

	// Cleanup when client disconnects
	log.Println("Client disconnected")
	mu.Lock()
	delete(clients, c)
	delete(typingUsers, c)
	mu.Unlock()
}
