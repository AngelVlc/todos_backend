package subscribers

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type UpdateSearchIndexDocumentProcessor struct {
	eventName         string
	eventBus          events.EventBus
	channel           chan events.DataEvent
	listsRepo         domain.ListsRepository
	listsSearchClient search.SearchIndexClient
	doneFunc          func(listID int32, err error)
	newRelicApp       *newrelic.Application
}

func NewUpdateSearchIndexDocumentProcessor(eventName string, eventBus events.EventBus, listsRepo domain.ListsRepository, listsSearchClient search.SearchIndexClient, newRelicApp *newrelic.Application) *UpdateSearchIndexDocumentProcessor {
	doneFunc := func(listID int32, err error) {
		if err != nil {
			log.Printf("Creting or updating search index document for list with ID %v\n failed with error %v", listID, err)
			honeybadger.Notify(err)
		}
	}

	return &UpdateSearchIndexDocumentProcessor{
		eventName:         eventName,
		eventBus:          eventBus,
		channel:           make(chan events.DataEvent),
		listsRepo:         listsRepo,
		listsSearchClient: listsSearchClient,
		doneFunc:          doneFunc,
		newRelicApp:       newRelicApp,
	}
}

func (s *UpdateSearchIndexDocumentProcessor) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *UpdateSearchIndexDocumentProcessor) Start() {
	for d := range s.channel {
		listID, _ := d.Data.(int32)
		txn := s.newRelicApp.StartTransaction(s.eventName)
		ctx := newrelic.NewContext(context.Background(), txn)

		srv := application.NewAddListToSearchIndexService(s.listsRepo, s.listsSearchClient)
		err := srv.AddListToSearchIndexService(ctx, listID)

		s.doneFunc(listID, err)

		txn.End()
	}
}
