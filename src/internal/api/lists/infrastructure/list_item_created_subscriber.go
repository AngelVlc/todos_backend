package infrastructure

import (
	"log"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemCreatedEventSubscriber struct {
	topic     string
	eventBus  events.EventBus
	channel   chan events.DataEvent
	listsRepo domain.ListsRepository
}

func NewListItemCreatedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository) *ListItemCreatedEventSubscriber {
	return &ListItemCreatedEventSubscriber{
		topic:     "listItemCreated",
		eventBus:  eventBus,
		channel:   make(chan events.DataEvent),
		listsRepo: listsRepo,
	}
}

func (s *ListItemCreatedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.topic, s.channel)
}

func (s *ListItemCreatedEventSubscriber) Start() {
	for {
		select {
		case d := <-s.channel:
			listID, _ := d.Data.(int32)
			log.Printf("Incrementing items counter for list with ID %v\n", d.Data)
			s.listsRepo.IncrementListCounter(listID)
		}
	}
}
