package comment

import (
	"time"

	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/eventBus"
	"gorm.io/gorm"
)

type CommentService struct {
	postgresDB  *gorm.DB
	postService *post.PostService
	eventBus    *eventbus.EventBus
}

func NewService(postgresDB *gorm.DB, postService *post.PostService, eventBus *eventbus.EventBus) *CommentService {
	return &CommentService{
		postgresDB:  postgresDB,
		postService: postService,
		eventBus:    eventBus,
	}
}

func (s *CommentService) CreateComment(req CommentRequest, postId uuid.UUID) (*Comment, error) {
	postIsPrivate, err := s.postService.IsPrivate(postId)
	if err != nil {
		return nil, err
	}

	comment, err := CreateComment(req, postId, postIsPrivate)
	if err != nil {
		return nil, err
	}

	result := s.postgresDB.Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}

	event := eventbus.Event{
		Type:      "UserRegistered",
		Timestamp: time.Now(),
		Data: eventbus.UserRegisteredEvent{
			ID:    1,
			Name:  "sdf",
			Email: "sdfs",
		},
	}

	s.eventBus.Publish(event)

	return comment, nil
}
