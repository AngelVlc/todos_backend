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

func TestIndexAllListsProcessor(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedSearchClient := search.MockedSearchIndexClient{}

	ctx := newrelic.NewContext(context.Background(), nil)

	foundLists := []domain.ListRecord{
		{ID: 11, UserID: 2, Name: "list1", Items: []*domain.ListItemRecord{{ID: 21, Title: "title1", Description: "desc1"}, {ID: 22, Title: "title2", Description: "desc2"}}},
		{ID: 12, UserID: 2, Name: "list2"},
	}
	mockedRepo.On("GetAllLists", ctx).Return(foundLists, nil).Once()

	l1vo, _ := domain.NewListNameValueObject("list1")
	l2vo, _ := domain.NewListNameValueObject("list2")

	listDocuments := []domain.ListSearchDocument{
		{ObjectID: "11", UserID: 2, Name: l1vo, ItemsTitles: []string{"title1", "title2"}, ItemsDescriptions: []string{"desc1", "desc2"}},
		{ObjectID: "12", UserID: 2, Name: l2vo, ItemsTitles: []string{}, ItemsDescriptions: []string{}},
	}
	mockedSearchClient.On("SaveObjects", listDocuments).Once().Return(nil)

	ch := make(chan events.DataEvent)
	doneChan := make(chan bool)
	f := func(err error) {
		doneChan <- true
	}
	subscriber := &IndexAllListsProcessor{
		channel:           ch,
		listsRepo:         &mockedRepo,
		listsSearchClient: &mockedSearchClient,
		doneFunc:          f,
	}

	go subscriber.Start()

	ch <- events.DataEvent{Data: nil}

	<-doneChan
	mockedRepo.AssertExpectations(t)
	mockedSearchClient.AssertExpectations(t)
}
