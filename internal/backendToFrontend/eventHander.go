package backendtofrontend

import (
	"fmt"

	"github.com/vladovsiychuk/microservice-demo-go/internal/comment"
	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/eventBus"
)

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
