package backendtofrontend

import (
	"fmt"

	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
)

type EventHandler struct {
	service BffServiceI
}

func NewEventHandler(service BffServiceI) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) PostCreatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		post, ok := event.Data.(*post.Post)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		h.service.CreatePostAggregate(post)
	}
}

func (h *EventHandler) PostUpdatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		post, ok := event.Data.(*post.Post)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		h.service.UpdatePostAggregate(post)
	}
}

func (h *EventHandler) CommentCreatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		comment, ok := event.Data.(*comment.Comment)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		h.service.AddCommentToPostAggregate(comment)
	}
}

func (h *EventHandler) CommentUpdatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		comment, ok := event.Data.(*comment.Comment)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		h.service.UpdateCommentInPostAggregate(comment)
	}
}
