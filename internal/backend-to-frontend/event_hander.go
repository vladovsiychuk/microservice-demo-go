package backendtofrontend

import (
	"fmt"

	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	"github.com/vladovsiychuk/microservice-demo-go/internal/post"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/event-bus"
)

func PostCreatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		post, ok := event.Data.(*post.Post)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New post created registered:")
		fmt.Println("content:", post.Content)
	}
}

func PostUpdatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		post, ok := event.Data.(*post.Post)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New post updated registered:")
		fmt.Println("content:", post.Content)
	}
}

func CommentCreatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		comment, ok := event.Data.(*comment.Comment)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New comment created registered:")
		fmt.Println("content:", comment.Content)
	}
}

func CommentUpdatedHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		comment, ok := event.Data.(*comment.Comment)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New comment updated registered:")
		fmt.Println("content:", comment.Content)
	}
}
