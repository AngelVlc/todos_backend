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
	l1vo, _ := domain.NewListNameValueObject("list1")
	l2vo, _ := domain.NewListNameValueObject("list2")
	t1vo, _ := domain.NewItemTitleValueObject("title1")
	d1vo, _ := domain.NewItemDescriptionValueObject("desc1")
	t2vo, _ := domain.NewItemTitleValueObject("title2")
	d2vo, _ := domain.NewItemDescriptionValueObject("desc2")

	foundLists := []*domain.ListEntity{
		{ID: 11, Name: l1vo, Items: []*domain.ListItemEntity{{ID: 21, Title: t1vo, Description: d1vo}, {ID: 22, Title: t2vo, Description: d2vo}}},
		{ID: 12, Name: l2vo},
	}
	mockedRepo.On("GetAllLists", ctx).Return(foundLists, nil).Once()

	listDocuments := []domain.ListSearchDocument{
		{ObjectID: "11", Name: l1vo, ItemsTitles: []string{"title1", "title2"}, ItemsDescriptions: []string{"desc1", "desc2"}},
		{ObjectID: "12", Name: l2vo, ItemsTitles: []string{}, ItemsDescriptions: []string{}},
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
