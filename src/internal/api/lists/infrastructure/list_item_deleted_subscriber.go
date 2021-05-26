package infrastructure

import (
	"log"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemDeletedEventSubscriber struct {
	topic     string
	eventBus  events.EventBus
	channel   chan events.DataEvent
	listsRepo domain.ListsRepository
}

func NewListItemDeletedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository) *ListItemDeletedEventSubscriber {
	return &ListItemDeletedEventSubscriber{
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
			listID, _ := d.Data.(int32)
			log.Printf("Decrementing items counter for list with ID %v\n", d.Data)
			s.listsRepo.DecrementListCounter(listID)
		}
	}
}
