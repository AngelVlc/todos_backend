package infrastructure

import (
	"context"
	"log"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ListItemDeletedEventSubscriber struct {
	eventName   string
	eventBus    events.EventBus
	channel     chan events.DataEvent
	listsRepo   domain.ListsRepository
	doneFunc    func(listID int32)
	newRelicApp *newrelic.Application
}

func NewListItemDeletedEventSubscriber(eventBus events.EventBus, listsRepo domain.ListsRepository, newRelicApp *newrelic.Application) *ListItemDeletedEventSubscriber {
	doneFunc := func(listID int32) {
		log.Printf("Decremented items counter for list with ID %v\n", listID)
	}

	return &ListItemDeletedEventSubscriber{
		eventName:   "listItemDeleted",
		eventBus:    eventBus,
		channel:     make(chan events.DataEvent),
		listsRepo:   listsRepo,
		doneFunc:    doneFunc,
		newRelicApp: newRelicApp,
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
			txn := s.newRelicApp.StartTransaction("decrementListCounter")
			ctx := newrelic.NewContext(context.Background(), txn)
			s.listsRepo.DecrementListCounter(ctx, listID)
			txn.End()
			s.doneFunc(listID)
		}
	}
}
