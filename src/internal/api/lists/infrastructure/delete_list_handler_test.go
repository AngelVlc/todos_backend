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
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func TestDeletesListHandler(t *testing.T) {
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

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the list fails", func(t *testing.T) {
		mockedRepo.On("FindListByID", int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := DeleteListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting the user list")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError if the list does not exits", func(t *testing.T) {
		mockedRepo.On("FindListByID", int32(11), int32(1)).Return(nil, nil).Once()

		result := DeleteListHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "The list does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the delete fails", func(t *testing.T) {
		list := domain.List{ID: 11, Name: "list1"}
		mockedRepo.On("FindListByID", int32(11), int32(1)).Return(&list, nil).Once()
		mockedRepo.On("DeleteList", int32(11), int32(1)).Return(fmt.Errorf("some error")).Once()

		result := DeleteListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error deleting the user list")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should delete the list", func(t *testing.T) {
		list := domain.List{ID: 11, Name: "list1"}
		mockedRepo.On("FindListByID", int32(11), int32(1)).Return(&list, nil).Once()
		mockedRepo.On("DeleteList", int32(11), int32(1)).Return(nil).Once()

		result := DeleteListHandler(httptest.NewRecorder(), request(), h)

		results.CheckOkResult(t, result, http.StatusNoContent)
		mockedRepo.AssertExpectations(t)
	})
}
