package websockets

import (
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	clients        map[*Client]bool
	broadcast      chan *Message
	broadactExcept chan *BroadcastExceptReq
	register       chan *Client
	unregister     chan *Client
}

type BroadcastExceptReq struct {
	message        *Message
	clientToIgnore uuid.UUID
}

func NewHub() *Hub {
	hub := &Hub{
		clients:        make(map[*Client]bool),
		broadcast:      make(chan *Message),
		broadactExcept: make(chan *BroadcastExceptReq),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}

	go hub.run()
	return hub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					log.Println("Client send buffer is full.")
					close(client.send)
					delete(h.clients, client)
				}
			}
		case req := <-h.broadactExcept:
			for client := range h.clients {
				if client.Id == req.clientToIgnore {
					continue
				}

				select {
				case client.send <- req.message:
				default:
					log.Println("Client send buffer is full.")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
