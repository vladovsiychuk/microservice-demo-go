package comment

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vladovsiychuk/microservice-demo-go/internal/shared"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
)

type CommentService struct {
	repository  CommentRepositoryI
	postService shared.PostServiceSharedI
	eventBus    eventbus.EventBusI
}

func NewService(
	repository CommentRepositoryI,
	postService shared.PostServiceSharedI,
	eventBus eventbus.EventBusI,
) *CommentService {
	return &CommentService{
		repository:  repository,
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

	if err := s.repository.Create(comment); err != nil {
		return nil, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:      shared.CommentCreatedEventType,
		Timestamp: time.Now(),
		Data:      comment,
	})

	return comment, nil
}

func (s *CommentService) UpdateComment(req CommentRequest, commentId uuid.UUID) (*Comment, error) {
	var comment Comment

	if err := s.repository.FindByKey(&comment, commentId); err != nil {
		return nil, errors.New("Comment not found")
	}

	if err := comment.Update(req); err != nil {
		return nil, err
	}

	if err := s.repository.Update(&comment); err != nil {
		return nil, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:      shared.CommentUpdatedEventType,
		Timestamp: time.Now(),
		Data:      &comment,
	})

	return &comment, nil
}
