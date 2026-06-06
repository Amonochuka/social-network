package chats

import (
	"log"
	"net/http"
	"time"
	"socialnetwork/dummydb"
	"github.com/gorilla/websocket"
	"socialnetwork/handlers"
)

func generateID() string {
    return time.Now().Format("20060102150405.000000")
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool   // registered clients
	broadcast  chan []byte        // inbound messages from clients
	register   chan *Client       // register requests from clients
	unregister chan *Client       // unregister requests from clients
	
}

// NewHub creates a new Hub instance





func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's main loop — processes register, unregister, and broadcast
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered. Total: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered. Total: %d", len(h.clients))
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

// Client represents a single WebSocket connection
type Client struct {
    hub       *Hub
    conn      *websocket.Conn
    send      chan []byte

    ID        string
    GroupID   string   // ADD THIS
    Email     string
    FirstName string
    LastName  string
}


type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

func GetUserByID(userID string) (*User, error) {
	var user User
	err := dummydb.DB.QueryRow(
		"SELECT id, email, first_name, last_name FROM users WHERE id = ?", 
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}





// readPump reads messages from the WebSocket connection to the hub
func (c *Client) readPump() {

	db := dummydb.DB;

	stmt, err := db.Prepare(`
    INSERT INTO group_messages (id, group_id, sender_id, content)
    VALUES (?, ?, ?, ?)
`)
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()



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
		////////INSERT SYNTAX ////////////
		    // INSERT INTO SQLITE
    _, err = stmt.Exec(
        generateID(),      // message id (you must implement)
        c.GroupID,         // group id
        c.ID,          // sender id
        string(message),   // content
    )

    if err != nil {
        log.Printf("DB insert error: %v", err)
    }

		
	}
}

// writePump writes messages from the hub to the WebSocket connection
func (c *Client) writePump() {
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

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

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

// ServeWs handles WebSocket requests from clients
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := handlers.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
		ID:        "anonymous",
		Email:     "anon@example.com",
		FirstName: "Anonymous",
		LastName:  "User",
	}

	// Optional: try to get user from DB if ID provided
	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		if dbClient, err := GetUserByID(userID); err == nil {
			client.ID = dbClient.ID
			client.Email = dbClient.Email
			client.FirstName = dbClient.FirstName
			client.LastName = dbClient.LastName
		}
	}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}
