package infrastructure

import (
	"log"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

type ListItemCreatedEventSubscriber struct {
	eventName string
	eventBus  events.EventBus
	channel   chan events.DataEvent
	listsRepo domain.ListsRepository
	doneFunc  func(listID int32)
}

func NewListItemCreatedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository) *ListItemCreatedEventSubscriber {
	doneFunc := func(listID int32) {
		log.Printf("Increment items counter for list with ID %v\n", listID)
	}

	return &ListItemCreatedEventSubscriber{
		eventName: "listItemCreated",
		eventBus:  eventBus,
		channel:   make(chan events.DataEvent),
		listsRepo: listsRepo,
		doneFunc:  doneFunc,
	}
}

func (s *ListItemCreatedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *ListItemCreatedEventSubscriber) Start() {
	for {
		select {
		case d := <-s.channel:
			listID, _ := d.Data.(int32)
			s.listsRepo.IncrementListCounter(listID)
			s.doneFunc(listID)
		}
	}
}
