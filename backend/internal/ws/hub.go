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

// broadcastPresence notifies every connected client that a user's
// online status has changed.
func (h *Hub) broadcastPresence(userID string, online bool) {
	payload := map[string]interface{}{
		"type":      "presence",
		"user_id":   userID,
		"is_online": online,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, clients := range h.clients {
		for _, client := range clients {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// Register adds a client connection for a user.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()

	// Remember whether this user was previously offline.
	wasOffline := len(h.clients[client.UserID]) == 0

	h.clients[client.UserID] = append(h.clients[client.UserID], client)

	// Send the newly connected client everyone who is already online.
	onlineUsers := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		onlineUsers = append(onlineUsers, userID)
	}

	payload := map[string]interface{}{
		"type":         "presence_update",
		"online_users": onlineUsers,
	}

	data, _ := json.Marshal(payload)

	h.mu.Unlock()

	client.Send <- data

	// Only broadcast when this is the user's first open tab.
	if wasOffline {
		h.broadcastPresence(client.UserID, true)
	}

	log.Printf("ws: user %s connected (total connections: %d)", client.UserID, len(h.clients[client.UserID]))
}

// Unregister removes a client connection when it disconnects.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()

	conns := h.clients[client.UserID]
	for i, c := range conns {
		if c == client {
			h.clients[client.UserID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// Only announce offline when every tab has disconnected.
	nowOffline := len(h.clients[client.UserID]) == 0

	if nowOffline {
		delete(h.clients, client.UserID)
	}

	h.mu.Unlock()

	close(client.Send)

	if nowOffline {
		h.broadcastPresence(client.UserID, false)
	}

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

// SendToGroup delivers a message to every connected member of a group.
func (h *Hub) SendToGroup(memberIDs []string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("ws: failed to marshal group payload:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range memberIDs {
		for _, client := range h.clients[userID] {
			select {
			case client.Send <- data:
			default:
				log.Printf("ws: send buffer full for user %s in group broadcast, dropping", userID)
			}
		}
	}
}