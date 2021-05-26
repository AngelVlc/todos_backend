//+build !e2e

package infrastructure

import (
	"testing"

	listsRepository "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
)

func TestListItemDeletedEventSubscriber(t *testing.T) {
	ch := make(chan events.DataEvent)
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedRepo.On("DecrementListCounter", int32(11)).Return(nil).Once()

	doneChan := make(chan bool)
	f := func(listID int32) {
		doneChan <- true
	}

	subscriber := &ListItemDeletedEventSubscriber{
		channel:   ch,
		listsRepo: &mockedRepo,
		doneFunc:  f,
	}

	go subscriber.Start()

	ch <- events.DataEvent{Data: int32(11)}

	<-doneChan
	mockedRepo.AssertExpectations(t)
}
