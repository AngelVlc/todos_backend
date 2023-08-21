package subscribers

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ListItemsCountProcessor struct {
	eventName   string
	eventBus    events.EventBus
	channel     chan events.DataEvent
	listsRepo   domain.ListsRepository
	doneFunc    func(listID int32, err error)
	newRelicApp *newrelic.Application
}

func NewListItemsCountProcessor(eventName string, eventBus events.EventBus, listsRepo domain.ListsRepository, newRelicApp *newrelic.Application) *ListItemsCountProcessor {
	doneFunc := func(listID int32, err error) {
		if err != nil {
			log.Printf("Updated items counter for list with ID %v\n failed with error %v", listID, err)
			honeybadger.Notify(err)
		} else {
			log.Printf("Updated items counter for list with ID %v\n", listID)
		}
	}

	return &ListItemsCountProcessor{
		eventName:   eventName,
		eventBus:    eventBus,
		channel:     make(chan events.DataEvent),
		listsRepo:   listsRepo,
		doneFunc:    doneFunc,
		newRelicApp: newRelicApp,
	}
}

func (s *ListItemsCountProcessor) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *ListItemsCountProcessor) Start() {
	for d := range s.channel {
		listID, _ := d.Data.(int32)
		txn := s.newRelicApp.StartTransaction("listItemsCountProcessor")
		ctx := newrelic.NewContext(context.Background(), txn)

		srv := application.NewUpdateListItemsCountService(s.listsRepo)
		err := srv.UpdateListsItemsCount(ctx, listID)
		s.doneFunc(listID, err)

		txn.End()
	}
}
