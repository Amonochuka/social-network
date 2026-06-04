package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	From    string `json:"from"`
	Content string `json:"content"`
}

type PrivateClient struct {
	hub      *PrivateHub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

type PrivateHub struct {
	clients       map[*PrivateClient]bool
	clientsByName map[string]*PrivateClient
	broadcast     chan []byte
	register      chan *PrivateClient
	unregister    chan *PrivateClient
}

func NewPrivateHub() *PrivateHub {
	return &PrivateHub{
		broadcast:     make(chan []byte),
		register:      make(chan *PrivateClient),
		unregister:    make(chan *PrivateClient),
		clients:       make(map[*PrivateClient]bool),
		clientsByName: make(map[string]*PrivateClient),
	}
}

func (h *PrivateHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientsByName[client.username] = client
			log.Printf("%s joined. Total: %d", client.username, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByName, client.username)
				close(client.send)
				log.Printf("%s left. Total: %d", client.username, len(h.clients))
			}

		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Invalid message:", err)
				continue
			}

			if msg.Type == "private" && msg.To != "" {
				if recipient, ok := h.clientsByName[msg.To]; ok {
					recipient.send <- message
				} else {
					log.Printf("User %s not found", msg.To)
				}
			} else {
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
						delete(h.clientsByName, client.username)
					}
				}
			}
		}
	}
}

func (c *PrivateClient) readPump() {
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

func (c *PrivateClient) writePump() {
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

func ServePrivateWs(hub *PrivateHub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &PrivateClient{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}