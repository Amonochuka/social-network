package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
)

// Upgrader configures how HTTP connections are migrated to WebSockets.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for local development. Modify for production security.
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	outbound chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:     conn,
		outbound: make(chan []byte, 256),
	}
}

// ServeWS handles incoming HTTP requests initiating WebSocket upgrade protocol.
func ServeWS(
	w http.ResponseWriter,
	r *http.Request,
	onConnect func(c *Client),
	onDisconnect func(),
	onMessage func(msg []byte),
) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(conn)

	// Trigger the hook to register this client inside our Hub maps
	onConnect(client)

	// Spin up concurrent loops for managing data traffic.
	go client.WritePump()

	// Pass the message and unregister hooks into the processing loop
	go client.ReadPump(onDisconnect, onMessage)
}

func (c *Client) WritePumpOnce(msg []byte) {
	c.send(msg)
}

func (c *Client) send(msg []byte) {
	select {
	case c.outbound <- msg:
	default:
		c.conn.Close()
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.outbound:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump(onDisconnect func(), msgHandler func([]byte)) {
	defer func() {
		onDisconnect() // Trigger the hook to remove client from Hub maps on exit
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("ws read error: %v", err)
			}
			return
		}
		msgHandler(msg)
	}
}
