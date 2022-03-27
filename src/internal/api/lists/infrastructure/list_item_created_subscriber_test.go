//go:build !e2e
// +build !e2e

package infrastructure

import (
	"context"
	"testing"

	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func TestListItemCreatedEventSubscriber(t *testing.T) {
	ch := make(chan events.DataEvent)
	mockedRepo := listsRepository.MockedListsRepository{}
	ctx := newrelic.NewContext(context.Background(), nil)
	mockedRepo.On("IncrementListCounter", ctx, int32(11)).Return(nil).Once()

	doneChan := make(chan bool)
	f := func(listID int32) {
		doneChan <- true
	}

	subscriber := &ListItemCreatedEventSubscriber{
		channel:   ch,
		listsRepo: &mockedRepo,
		doneFunc:  f,
	}

	go subscriber.Start()

	ch <- events.DataEvent{Data: int32(11)}

	<-doneChan
	mockedRepo.AssertExpectations(t)
}
