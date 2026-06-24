package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second    // max time to write a message to the peer
	pongWait   = 60 * time.Second    // max time to wait for a pong from the peer
	pingPeriod = (pongWait * 9) / 10 // send ping slightly before pong timeout
	maxMsgSize = 4096                // max incoming message size in bytes
)

// Client wraps a single gorilla WebSocket connection.
// It has two goroutines: readPump (blocks waiting for incoming messages)
// and writePump (blocks waiting for messages to send out).
// They communicate via the outbound channel.
type Client struct {
	conn     *websocket.Conn
	outbound chan []byte // hub writes here; writePump reads from here
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:     conn,
		outbound: make(chan []byte, 256), // buffered so hub never blocks
	}
}

// WritePumpOnce enqueues a single message immediately (used to push history on connect).
func (c *Client) WritePumpOnce(msg []byte) {
	c.send(msg)
}

// send is called by the hub to enqueue an outbound message.
// Non-blocking: if the buffer is full the client is considered too slow and dropped.
func (c *Client) send(msg []byte) {
	select {
	case c.outbound <- msg:
	default:
		// client is too slow; close its connection
		c.conn.Close()
	}
}

// WritePump reads from outbound and writes to the WebSocket.
// One goroutine per client. Must be run as `go client.WritePump()`.
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
				// hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			// send a ping to keep the connection alive
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump reads incoming messages from the WebSocket and passes them to msgHandler.
// One goroutine per client. Must be run as `go client.ReadPump(...)`.
// When the connection closes (browser tab closes, network drop, etc.) ReadPump returns,
// and the caller's defer runs cleanup (unregister from hub).
func (c *Client) ReadPump(msgHandler func([]byte)) {
	defer c.conn.Close()

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
			return // exits loop → deferred conn.Close runs
		}
		msgHandler(msg)
	}
}
