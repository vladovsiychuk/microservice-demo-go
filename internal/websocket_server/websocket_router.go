package websocketserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketRouter struct{}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewRouter() *WebSocketRouter {
	return &WebSocketRouter{}
}

func (w *WebSocketRouter) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws", w.handleWebSocket)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (w *WebSocketRouter) handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}
	defer ws.Close()

	for {
		var message Message

		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		response := Message{
			Type:    "response",
			Content: fmt.Sprintf("Received: %s", message.Content),
		}
		if err := ws.WriteJSON(response); err != nil {
			log.Println("Error writing JSON response:", err)
			break
		}
	}
}
