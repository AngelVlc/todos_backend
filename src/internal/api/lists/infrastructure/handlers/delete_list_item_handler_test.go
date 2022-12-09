//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func TestDeletesListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_ListItem_Fails(t *testing.T) {
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

	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemEntity{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListItemHandler_Returns_An_Error_With_An_UnexpectedError_If_The_Delete_Fails(t *testing.T) {
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

	listItem := domain.ListItemEntity{ID: 111, ListID: 11, Title: "title"}
	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemEntity{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(&listItem, nil).Once()
	mockedRepo.On("DeleteListItem", request().Context(), int32(111), int32(11), int32(1)).Return(fmt.Errorf("some error")).Once()

	result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the list item")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListItemHandler_Deleted_The_ListItem(t *testing.T) {
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

	listItem := domain.ListItemEntity{ID: 111, ListID: 11, Title: "title"}
	mockedRepo.On("FindListItem", request().Context(), &domain.ListItemEntity{ID: int32(111), ListID: int32(11), UserID: int32(1)}).Return(&listItem, nil).Once()
	mockedRepo.On("DeleteListItem", request().Context(), int32(111), int32(11), int32(1)).Return(nil).Once()
	mockedEventBus.On("Publish", "listItemDeleted", int32(11))

	mockedEventBus.Wg.Add(1)
	result := DeleteListItemHandler(httptest.NewRecorder(), request(), h)
	mockedEventBus.Wg.Wait()

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
