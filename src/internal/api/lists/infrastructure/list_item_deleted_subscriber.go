package infrastructure

import (
	"log"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemDeletedEventSubscriber struct {
	eventName string
	eventBus  events.EventBus
	channel   chan events.DataEvent
	listsRepo domain.ListsRepository
	doneFunc  func(listID int32)
}

func NewListItemDeletedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository) *ListItemDeletedEventSubscriber {
	doneFunc := func(listID int32) {
		log.Printf("Decremented items counter for list with ID %v\n", listID)
	}

	return &ListItemDeletedEventSubscriber{
		eventName: "listItemDeleted",
		eventBus:  eventBus,
		channel:   make(chan events.DataEvent),
		listsRepo: listsRepo,
		doneFunc:  doneFunc,
	}
}

func (s *ListItemDeletedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *ListItemDeletedEventSubscriber) Start() {
	for {
		select {
		case d := <-s.channel:
			listID, _ := d.Data.(int32)
			s.listsRepo.DecrementListCounter(listID)
			s.doneFunc(listID)
		}
	}
}
