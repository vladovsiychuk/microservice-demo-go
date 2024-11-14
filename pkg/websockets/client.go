package websockets

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 10 * time.Second
)

type Client struct {
	Id   uuid.UUID
	hub  *Hub
	conn *websocket.Conn
	send chan *Message
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func CreateClient(hub *Hub, id uuid.UUID, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}

	client := &Client{
		Id:   id,
		hub:  hub,
		conn: ws,
		send: make(chan *Message, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var message Message

		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadactExcept <- &BroadcastExceptReq{
			message:        &message,
			clientToIgnore: c.Id,
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteJSON(websocket.CloseMessage)
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Println("Error writing JSON response:", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteJSON(websocket.PingMessage); err != nil {
				return
			}
		}
	}
}
