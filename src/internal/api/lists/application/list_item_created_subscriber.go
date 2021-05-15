package application

import (
	"fmt"
	"time"

	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemCreatedEventSubscriber struct {
	topic    string
	eventBus events.EventBus
	channel  chan events.DataEvent
}

func NewListItemCreatedEventSubscriber(eventBus events.EventBus) *ListItemCreatedEventSubscriber {
	return &ListItemCreatedEventSubscriber{
		topic:    "listItemCreated",
		eventBus: eventBus,
		channel:  make(chan events.DataEvent),
	}
}

func (s *ListItemCreatedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.topic, s.channel)
}

func (s *ListItemCreatedEventSubscriber) Start() {

	for {
		select {
		case d := <-s.channel:
			// go printDataEvent("ch1", d)
			printDataEvent("ch1", d)
		}
	}
}

func printDataEvent(ch string, data events.DataEvent) {
	time.Sleep(3 * time.Second)
	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}
