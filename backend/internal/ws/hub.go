package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents one connected user's websocket connection.
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub keeps track of all connected clients, keyed by userID.
// A user could theoretically have multiple tabs open, so we store a slice per user.
type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string][]*Client),
	}
}

// Register adds a client connection for a user.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.UserID] = append(h.clients[client.UserID], client)
	log.Printf("ws: user %s connected (total connections: %d)", client.UserID, len(h.clients[client.UserID]))
}

// Unregister removes a client connection when it disconnects.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.clients[client.UserID]
	for i, c := range conns {
		if c == client {
			h.clients[client.UserID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.clients[client.UserID]) == 0 {
		delete(h.clients, client.UserID)
	}
	close(client.Send)
	log.Printf("ws: user %s disconnected", client.UserID)
}

// SendToUser delivers a message to every active connection for a given user.
// If the user isn't connected, the message is simply dropped (they'll see it
// next time they fetch message history via the REST endpoint).
func (h *Hub) SendToUser(userID string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("ws: failed to marshal payload:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients[userID] {
		select {
		case client.Send <- data:
		default:
			log.Printf("ws: send buffer full for user %s, dropping message", userID)
		}
	}
}

// IsOnline reports whether a user currently has at least one active connection.
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}
