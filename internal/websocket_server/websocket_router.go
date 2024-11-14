package websocketserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/pkg/websockets"
)

type WebSocketRouter struct {
	service *WsService
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewRouter(service *WsService) *WebSocketRouter {
	return &WebSocketRouter{
		service: service,
	}
}

func (w *WebSocketRouter) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws/rooms/:roomId/users/:userId", w.handleWebSocket)
}

func (w *WebSocketRouter) handleWebSocket(c *gin.Context) {
	roomIdStr := c.Param("roomId")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	userIdStr := c.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	room := w.service.GetOrCreateRoom(roomId)
	websockets.CreateClient(room, userId, c.Writer, c.Request)
}
