package eventbus

import (
	"time"
)

type Event struct {
	Type      string
	Timestamp time.Time
	Data      interface{}
}

type UserRegisteredEvent struct {
	ID    int
	Name  string
	Email string
}

type EventBus struct {
	subscribers map[string][]chan<- Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
	}
}

func (eb *EventBus) Subscribe(eventType string, subscriber chan<- Event) {
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

func (eb *EventBus) Publish(event Event) {
	subscribers := eb.subscribers[event.Type]
	for _, subscriber := range subscribers {
		subscriber <- event
	}
}
