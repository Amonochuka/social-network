package ws

import "sync"

// Hub is the single in-memory registry of all live WebSocket connections.
// Every connected client registers itself here; every outbound message is
// delivered through here. There is ONE hub for the whole server.
//
// Think of it like a switchboard operator: clients plug in, and when a
// message arrives for user X, the hub looks up X's socket and writes to it.
type Hub struct {
	mu sync.RWMutex // protects the maps below from concurrent reads/writes

	// privateClients: userID → Client
	// One user can only have one active private-chat connection at a time.
	privateClients map[string]*Client

	// groupClients: groupID → set of Clients
	// Multiple users are in the same group, so the value is a map (used as a set).
	groupClients map[string]map[*Client]bool

	// notifClients: userID → Client
	// Same idea as privateClients but for the notifications socket.
	notifClients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		privateClients: make(map[string]*Client),
		groupClients:   make(map[string]map[*Client]bool),
		notifClients:   make(map[string]*Client),
	}
}

// ── Registration ─────────────────────────────────────────────────────────────

func (h *Hub) RegisterPrivate(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.privateClients[userID] = c
}

func (h *Hub) UnregisterPrivate(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.privateClients, userID)
}

func (h *Hub) RegisterGroup(groupID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.groupClients[groupID] == nil {
		h.groupClients[groupID] = make(map[*Client]bool)
	}
	h.groupClients[groupID][c] = true
}

func (h *Hub) UnregisterGroup(groupID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if members, ok := h.groupClients[groupID]; ok {
		delete(members, c)
		if len(members) == 0 {
			delete(h.groupClients, groupID)
		}
	}
}

func (h *Hub) RegisterNotif(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notifClients[userID] = c
}

func (h *Hub) UnregisterNotif(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.notifClients, userID)
}

// ── Delivery ─────────────────────────────────────────────────────────────────

// SendToUser delivers a JSON payload to a user's private-chat socket.
// If the user is offline, the message is silently dropped (it was already
// persisted to SQLite, so nothing is lost).
func (h *Hub) SendToUser(userID string, payload []byte) {
	h.mu.RLock()
	c, ok := h.privateClients[userID]
	h.mu.RUnlock()
	if ok {
		c.send(payload)
	}
}

// BroadcastToGroup delivers a JSON payload to every client in a group safely.
// BroadcastToGroup delivers a JSON payload to every client in a group safely.
func (h *Hub) BroadcastToGroup(groupID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock() // Lock stays active until the loop finishes completely

	if members, ok := h.groupClients[groupID]; ok {
		for c := range members {
			c.send(payload)
		}
	}
}

// SendNotification delivers a notification to a user's notification socket.
func (h *Hub) SendNotification(userID string, payload []byte) {
	h.mu.RLock()
	c, ok := h.notifClients[userID]
	h.mu.RUnlock()
	if ok {
		c.send(payload)
	}
}
