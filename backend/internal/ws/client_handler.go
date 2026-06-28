package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production you'd restrict this to your actual frontend origin.
	// Since this is a local dev project behind our own CORS middleware already,
	// we allow it here too so the websocket handshake itself isn't blocked.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS upgrades the connection and starts the read/write pumps for one user.
// userID must already be authenticated by the caller (handler) before this runs.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws: upgrade failed:", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	hub.Register(client)

	go client.writePump()
	client.readPump(hub)
}

// readPump just keeps the connection alive and discards incoming pings/pongs.
// We don't expect the client to send arbitrary data over the socket — all
// actual message sending happens through the normal REST POST endpoint, which
// then pushes the result out over this socket via SendToUser. This keeps the
// permission checks (CanChat, group membership etc.) in one place: the service
// layer, instead of duplicating them inside the websocket handler.
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}