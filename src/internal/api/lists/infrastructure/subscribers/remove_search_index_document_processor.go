package subscribers

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type RemoveSearchIndexDocumentProcessor struct {
	eventName         string
	eventBus          events.EventBus
	channel           chan events.DataEvent
	listsSearchClient search.SearchIndexClient
	doneFunc          func(listID int32, err error)
	newRelicApp       *newrelic.Application
}

func NewRemoveSearchIndexDocumentProcessor(eventName string, eventBus events.EventBus, listsSearchClient search.SearchIndexClient, newRelicApp *newrelic.Application) *RemoveSearchIndexDocumentProcessor {
	doneFunc := func(listID int32, err error) {
		if err != nil {
			log.Printf("Removing search index document for list with ID %v\n failed with error %v", listID, err)
		}
	}

	return &RemoveSearchIndexDocumentProcessor{
		eventName:         eventName,
		eventBus:          eventBus,
		channel:           make(chan events.DataEvent),
		listsSearchClient: listsSearchClient,
		doneFunc:          doneFunc,
		newRelicApp:       newRelicApp,
	}
}

func (s *RemoveSearchIndexDocumentProcessor) Subscribe() {
	s.eventBus.Subscribe(s.eventName, s.channel)
}

func (s *RemoveSearchIndexDocumentProcessor) Start() {
	for {
		select {
		case d := <-s.channel:
			listID, _ := d.Data.(int32)
			txn := s.newRelicApp.StartTransaction(s.eventName)
			ctx := newrelic.NewContext(context.Background(), txn)

			srv := application.NewRemoveListFromSearchIndexService(s.listsSearchClient)
			err := srv.RemoveListFromSearchIndexService(ctx, listID)

			s.doneFunc(listID, err)

			txn.End()
		}
	}
}
