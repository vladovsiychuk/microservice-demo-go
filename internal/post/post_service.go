package post

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
	"gorm.io/gorm"
)

type PostService struct {
	postgresDB *gorm.DB
	eventBus   *eventbus.EventBus
}

func NewService(postgresDB *gorm.DB, eventBus *eventbus.EventBus) *PostService {
	return &PostService{
		postgresDB: postgresDB,
		eventBus:   eventBus,
	}
}

func (s *PostService) CreatePost(req PostRequest) (*Post, error) {
	post, err := CreatePost(req)
	if err != nil {
		return nil, err
	}

	result := s.postgresDB.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	s.eventBus.Publish(eventbus.Event{
		Type:      shared.PostCreatedEventType,
		Timestamp: time.Now(),
		Data:      post,
	})

	return post, nil
}

func (s *PostService) UpdatePost(postId uuid.UUID, req PostRequest) (*Post, error) {
	var post Post

	if err := s.postgresDB.Take(&post, postId).Error; err != nil {
		return nil, errors.New("Post not found")
	}

	if err := post.Update(req); err != nil {
		return nil, err
	}

	if err := s.postgresDB.Save(&post).Error; err != nil {
		return nil, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:      shared.PostUpdatedEventType,
		Timestamp: time.Now(),
		Data:      &post,
	})

	return &post, nil
}

func (s *PostService) IsPrivate(postId uuid.UUID) (bool, error) {
	var post Post

	if err := s.postgresDB.Take(&post, postId).Error; err != nil {
		return false, errors.New("Post not found")
	}

	return post.IsPrivate, nil
}
