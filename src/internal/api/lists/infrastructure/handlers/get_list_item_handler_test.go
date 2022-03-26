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

func TestGetListItemHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
			"id":     "111",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	t.Run("Should return an error if the query to find the list item fails", func(t *testing.T) {
		mockedRepo.On("FindListItemByID", request().Context(), int32(111), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return the list", func(t *testing.T) {
		list := domain.ListItem{ID: 111, ListID: 11, Title: "title", Description: "desc"}
		mockedRepo.On("FindListItemByID", request().Context(), int32(111), int32(11), int32(1)).Return(&list, nil).Once()

		result := GetListItemHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		listRes, isOk := okRes.Content.(*infrastructure.ListItemResponse)
		require.Equal(t, true, isOk, "should be a list item response")

		assert.Equal(t, int32(111), listRes.ID)
		assert.Equal(t, int32(11), listRes.ListID)
		assert.Equal(t, "title", listRes.Title)
		assert.Equal(t, "desc", listRes.Description)
		mockedRepo.AssertExpectations(t)
	})
}
