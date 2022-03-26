//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure"
	listsRepository "github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetListHandler(t *testing.T) {
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

	t.Run("Should return an error if the query to find the list fails", func(t *testing.T) {
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetListHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return the list", func(t *testing.T) {
		list := domain.List{ID: 11, Name: "list1", ItemsCount: 4}
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(&list, nil).Once()

		result := GetListHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		listRes, isOk := okRes.Content.(*infrastructure.ListResponse)
		require.Equal(t, true, isOk, "should be a list response")

		assert.Equal(t, int32(11), listRes.ID)
		assert.Equal(t, "list1", listRes.Name)
		assert.Equal(t, int32(4), listRes.ItemsCount)
		mockedRepo.AssertExpectations(t)
	})
}
