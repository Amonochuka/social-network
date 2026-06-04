package handlers

import (
	"log"
	"strings"
	"sync"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
)

// GroupHub manages one chat room
type GroupHub struct {
	clients    map[*GroupClient]bool
	broadcast  chan []byte
	register   chan *GroupClient
	unregister chan *GroupClient
	roomName   string
}

// NewGroupHub creates a hub for one room
func NewGroupHub(roomName string) *GroupHub {
	return &GroupHub{
		broadcast:  make(chan []byte),
		register:   make(chan *GroupClient),
		unregister: make(chan *GroupClient),
		clients:    make(map[*GroupClient]bool),
		roomName:   roomName,
	}
}

// GroupHub Run method
func (h *GroupHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("%s joined %s. Total: %d", client.username, h.roomName, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("%s left %s. Total: %d", client.username, h.roomName, len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// GroupClient struct
type GroupClient struct {
	hub      *GroupHub
	conn     *websocket.Conn
	send     chan []byte
	username string
	room     string
}

// RoomManager
type RoomManager struct {
	rooms map[string]*GroupHub
	mu    sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*GroupHub),
	}
}

func (rm *RoomManager) GetOrCreateRoom(name string) *GroupHub {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[name]; exists {
		return room
	}

	log.Printf("Creating new room: %s", name)
	room := NewGroupHub(name)
	rm.rooms[name] = room
	go room.Run()
	return room
}

func ExtractRoomName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}


func ServeGroupWs(hub *GroupHub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &GroupClient{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		room:     hub.roomName,
	}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}


func (c *GroupClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *GroupClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}