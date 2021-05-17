package infrastructure

import (
	"fmt"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemDeletedEventSubscriber struct {
	topic     string
	eventBus  events.EventBus
	channel   chan events.DataEvent
	listsRepo domain.ListsRepository
}

func NewListItemDeletedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository) *ListItemCreatedEventSubscriber {
	return &ListItemCreatedEventSubscriber{
		topic:     "listItemDeleted",
		eventBus:  eventBus,
		channel:   make(chan events.DataEvent),
		listsRepo: listsRepo,
	}
}

func (s *ListItemDeletedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.topic, s.channel)
}

func (s *ListItemDeletedEventSubscriber) Start() {
	for {
		select {
		case d := <-s.channel:
			// go printDataEvent("ch1", d)
			fmt.Printf("Topic: %s; DataEvent: %v\n", d.Topic, d.Data)
		}
	}
}
