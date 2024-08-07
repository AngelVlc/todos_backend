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

func deleteRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestDeletesListHandler_Returns_An_Error_If_The_Query_To_Find_The_Existing_List_Fails(t *testing.T) {
	request := deleteRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := DeleteListHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Delete_Fails(t *testing.T) {
	request := deleteRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	existingList := domain.ListRecord{ID: 11, Name: "list1"}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&existingList, nil).Once()
	mockedRepo.On("DeleteList", request.Context(), existingList).Return(fmt.Errorf("some error")).Once()

	result := DeleteListHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the user list")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListHandler_Deletes_The_List(t *testing.T) {
	request := deleteRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
	}

	existingList := domain.ListRecord{ID: 11, Name: "list1", UserID: 1}
	mockedRepo.On("FindList", request.Context(), domain.ListRecord{ID: 11, UserID: 1}).Return(&existingList, nil).Once()
	mockedRepo.On("DeleteList", request.Context(), existingList).Return(nil).Once()

	mockedEventBus.On("Publish", events.ListDeleted, int32(11))

	mockedEventBus.Wg.Add(1)
	result := DeleteListHandler(httptest.NewRecorder(), request, h)
	mockedEventBus.Wg.Wait()

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
