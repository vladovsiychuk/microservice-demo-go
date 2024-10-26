package backendtofrontend

import (
	"fmt"

	eventbus "github.com/vladovsiychuk/microservice-demo-go/pkg/eventBus"
)

func UserRegisteredHandler(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		userRegisteredEvent, ok := event.Data.(eventbus.UserRegisteredEvent)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New user registered:")
		fmt.Println("ID:", userRegisteredEvent.ID)
		fmt.Println("Name:", userRegisteredEvent.Name)
		fmt.Println("Email:", userRegisteredEvent.Email)
	}
}

func UserRegisteredHandler2(eventChan <-chan eventbus.Event) {
	for event := range eventChan {
		userRegisteredEvent, ok := event.Data.(eventbus.UserRegisteredEvent)
		if !ok {
			fmt.Println("Invalid event data")
			continue
		}

		// Handle the event
		fmt.Println("New user registered:")
		fmt.Println("ID:", userRegisteredEvent.ID)
		fmt.Println("Name:", userRegisteredEvent.Name)
		fmt.Println("Email:", userRegisteredEvent.Email)
	}
}
