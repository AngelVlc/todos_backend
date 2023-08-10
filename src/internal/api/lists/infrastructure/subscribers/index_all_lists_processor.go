package subscribers

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type IndexAllListsProcessor struct {
	eventName         string
	eventBus          events.EventBus
	channel           chan events.DataEvent
	listsRepo         domain.ListsRepository
	listsSearchClient search.SearchIndexClient
	doneFunc          func(err error)
	newRelicApp       *newrelic.Application
}

func NewIndexAllListsProcessor(eventName string, eventBus events.EventBus, listsRepo domain.ListsRepository, listsSearchClient search.SearchIndexClient, newRelicApp *newrelic.Application) *IndexAllListsProcessor {
	doneFunc := func(err error) {
		if err != nil {
			log.Printf("Index all lists failed with error %v", err)
		}
	}

	return &IndexAllListsProcessor{
		eventName:         eventName,
		eventBus:          eventBus,
		channel:           make(chan events.DataEvent),
		listsRepo:         listsRepo,
		listsSearchClient: listsSearchClient,
		doneFunc:          doneFunc,
		newRelicApp:       newRelicApp,
	}
}

func (s *IndexAllListsProcessor) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *IndexAllListsProcessor) Start() {
	for {
		select {
		case <-s.channel:
			txn := s.newRelicApp.StartTransaction(s.eventName)
			ctx := newrelic.NewContext(context.Background(), txn)

			srv := application.NewIndexAllListsService(s.listsRepo, s.listsSearchClient)
			err := srv.IndexAllLists(ctx)

			s.doneFunc(err)

			txn.End()
		}
	}
}
