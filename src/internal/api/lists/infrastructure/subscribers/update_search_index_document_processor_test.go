package subscribers

import (
	"context"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func TestUpdateSearchIndexDocumentProcessor(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedSearchClient := search.MockedSearchIndexClient{}

	ctx := newrelic.NewContext(context.Background(), nil)
	l1vo, _ := domain.NewListNameValueObject("list1")

	foundList := domain.ListEntity{
		ID:   12,
		Name: l1vo,
	}
	mockedRepo.On("FindList", ctx, domain.ListEntity{ID: 12}).Return(&foundList, nil).Once()

	listDocument := domain.ListSearchDocument{
		ObjectID:          "12",
		Name:              l1vo,
		ItemsTitles:       []string{},
		ItemsDescriptions: []string{},
	}
	mockedSearchClient.On("SaveObjects", listDocument).Once().Return(nil)

	ch := make(chan events.DataEvent)
	doneChan := make(chan bool)
	f := func(listID int32, err error) {
		doneChan <- true
	}
	subscriber := &UpdateSearchIndexDocumentProcessor{
		channel:           ch,
		listsRepo:         &mockedRepo,
		listsSearchClient: &mockedSearchClient,
		doneFunc:          f,
	}

	go subscriber.Start()

	ch <- events.DataEvent{Data: int32(12)}

	<-doneChan
	mockedRepo.AssertExpectations(t)
	mockedSearchClient.AssertExpectations(t)
}
