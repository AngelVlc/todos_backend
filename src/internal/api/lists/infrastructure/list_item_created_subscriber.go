package infrastructure

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/internal/api/shared/domain/events"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ListItemCreatedEventSubscriber struct {
	eventName   string
	eventBus    events.EventBus
	channel     chan events.DataEvent
	listsRepo   domain.ListsRepository
	doneFunc    func(listID int32)
	newRelicApp *newrelic.Application
}

func NewListItemCreatedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository, newRelicApp *newrelic.Application) *ListItemCreatedEventSubscriber {
	doneFunc := func(listID int32) {
		log.Printf("Incremented items counter for list with ID %v\n", listID)
	}

	return &ListItemCreatedEventSubscriber{
		eventName:   "listItemCreated",
		eventBus:    eventBus,
		channel:     make(chan events.DataEvent),
		listsRepo:   listsRepo,
		doneFunc:    doneFunc,
		newRelicApp: newRelicApp,
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
			txn := s.newRelicApp.StartTransaction("incrementListCounter")
			ctx := newrelic.NewContext(context.Background(), txn)
			s.listsRepo.IncrementListCounter(ctx, listID)
			txn.End()
			s.doneFunc(listID)
		}
	}
}
