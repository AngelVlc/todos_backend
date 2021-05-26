//+build !e2e

package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func TestDeletesListItemHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id":     "111",
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the list item fails", func(t *testing.T) {
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting the list item")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError if the list item does not exits", func(t *testing.T) {
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(nil, nil).Once()

		result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "The list item does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the delete fails", func(t *testing.T) {
		listItem := domain.ListItem{ID: 111, ListID: 11, Title: "title"}
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(&listItem, nil).Once()
		mockedRepo.On("DeleteListItem", int32(111), int32(11), int32(1)).Return(fmt.Errorf("some error")).Once()

		result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error deleting the list item")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should delete the list item", func(t *testing.T) {
		listItem := domain.ListItem{ID: 111, ListID: 11, Title: "title"}
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(&listItem, nil).Once()
		mockedRepo.On("DeleteListItem", int32(111), int32(11), int32(1)).Return(nil).Once()
		mockedEventBus.On("Publish", "listItemDeleted", int32(11))

		mockedEventBus.Wg.Add(1)
		result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)
		mockedEventBus.Wg.Wait()

		results.CheckOkResult(t, result, http.StatusNoContent)
		mockedRepo.AssertExpectations(t)
		mockedEventBus.AssertExpectations(t)
	})
}
