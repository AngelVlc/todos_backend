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
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func TestDeletesListHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	mockedRepo.On("FindList", request().Context(), &domain.ListEntity{ID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := DeleteListHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Delete_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	list := domain.ListEntity{ID: 11, Name: "list1"}
	mockedRepo.On("FindList", request().Context(), &domain.ListEntity{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("DeleteList", request().Context(), int32(11), int32(1)).Return(fmt.Errorf("some error")).Once()

	result := DeleteListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the user list")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesListHandler_Deletes_The_List(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	list := domain.ListEntity{ID: 11, Name: "list1"}
	mockedRepo.On("FindList", request().Context(), &domain.ListEntity{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("DeleteList", request().Context(), int32(11), int32(1)).Return(nil).Once()

	result := DeleteListHandler(httptest.NewRecorder(), request(), h)

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedRepo.AssertExpectations(t)
}
