package post

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
)

type PostService struct {
	repository PostRepositoryI
	eventBus   eventbus.EventBusI
}

type PostServiceI interface {
	CreatePost(req PostRequest) (*Post, error)
	UpdatePost(postId uuid.UUID, req PostRequest) (*Post, error)
}

func NewService(repository PostRepositoryI, eventBus eventbus.EventBusI) *PostService {
	return &PostService{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (s *PostService) CreatePost(req PostRequest) (*Post, error) {
	post, err := CreatePost(req)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(post); err != nil {
		return nil, err
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

	if err := s.repository.FindByKey(&post, postId); err != nil {
		return nil, errors.New("Post not found")
	}

	if err := post.Update(req); err != nil {
		return nil, err
	}

	if err := s.repository.Update(&post); err != nil {
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

	if err := s.repository.FindByKey(&post, postId); err != nil {
		return false, errors.New("Post not found")
	}

	return post.IsPrivate, nil
}
