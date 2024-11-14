package websocketserver

import (
	"github.com/google/uuid"
	ws "github.com/vladovsiychuk/microservice-demo-go/pkg/websockets"
)

type WsService struct {
	roomMap map[uuid.UUID]*ws.Hub

	getOrCreateRoom chan uuid.UUID
	roomResponse    chan *ws.Hub
}

func NewService() *WsService {
	service := &WsService{
		roomMap:         make(map[uuid.UUID]*ws.Hub),
		getOrCreateRoom: make(chan uuid.UUID),
		roomResponse:    make(chan *ws.Hub),
	}
	go service.run()
	return service
}

func (s *WsService) run() {
	for roomId := range s.getOrCreateRoom {
		hub, exists := s.roomMap[roomId]
		if !exists {
			hub = ws.NewHub()
			s.roomMap[roomId] = hub
		}
		s.roomResponse <- hub
	}
}

func (s *WsService) GetOrCreateRoom(roomId uuid.UUID) *ws.Hub {
	s.getOrCreateRoom <- roomId
	return <-s.roomResponse
}
