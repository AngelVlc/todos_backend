package subscribers

import (
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
)

func TestRemoveSearchIndexDocumentProcessor(t *testing.T) {
	mockedSearchClient := search.MockedSearchIndexClient{}

	mockedSearchClient.On("DeleteObject", "12").Once().Return(nil)

	ch := make(chan events.DataEvent)
	doneChan := make(chan bool)
	f := func(listID int32, err error) {
		doneChan <- true
	}
	subscriber := &RemoveSearchIndexDocumentProcessor{
		channel:           ch,
		listsSearchClient: &mockedSearchClient,
		doneFunc:          f,
	}

	go subscriber.Start()

	ch <- events.DataEvent{Data: int32(12)}

	<-doneChan
	mockedSearchClient.AssertExpectations(t)
}
