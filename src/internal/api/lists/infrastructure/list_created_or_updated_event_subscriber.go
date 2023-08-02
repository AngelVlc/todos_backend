package infrastructure

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ListCreatedOrUpdatedEventSubscriber struct {
	eventName   string
	eventBus    events.EventBus
	channel     chan events.DataEvent
	listsRepo   domain.ListsRepository
	doneFunc    func(listID int32)
	newRelicApp *newrelic.Application
}

func NewListCreatedOrUpdatedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository, newRelicApp *newrelic.Application) *ListCreatedOrUpdatedEventSubscriber {
	doneFunc := func(listID int32) {
		log.Printf("Updated items counter for list with ID %v\n", listID)
	}

	return &ListCreatedOrUpdatedEventSubscriber{
		eventName:   "listCreatedOrUpdated",
		eventBus:    eventBus,
		channel:     make(chan events.DataEvent),
		listsRepo:   listsRepo,
		doneFunc:    doneFunc,
		newRelicApp: newRelicApp,
	}
}

func (s *ListCreatedOrUpdatedEventSubscriber) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *ListCreatedOrUpdatedEventSubscriber) Start() {
	for {
		select {
		case d := <-s.channel:
			listID, _ := d.Data.(int32)
			txn := s.newRelicApp.StartTransaction("updateListCounter")
			ctx := newrelic.NewContext(context.Background(), txn)
			s.listsRepo.UpdateListItemsCounter(ctx, listID)
			txn.End()
			s.doneFunc(listID)
		}
	}
}
