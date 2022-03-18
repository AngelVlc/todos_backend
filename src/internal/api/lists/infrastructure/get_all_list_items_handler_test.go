//go:build !e2e
// +build !e2e

package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllListItemsHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	t.Run("Should return an error result with an unexpected error if the query fails", func(t *testing.T) {
		mockedRepo.On("GetAllListItems", request().Context(), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetAllListItemsHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting all list items")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return the lists if the query does not fail", func(t *testing.T) {
		found := []domain.ListItem{
			{ID: 111, ListID: 11, Title: "title1", Description: "desc1", Position: int32(0)},
			{ID: 112, ListID: 11, Title: "title2", Description: "desc2", Position: int32(1)},
		}

		mockedRepo.On("GetAllListItems", request().Context(), int32(11), int32(1)).Return(found, nil)

		result := GetAllListItemsHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		listRes, isOk := okRes.Content.([]ListItemResponse)
		require.Equal(t, true, isOk, "should be an array of list item response")

		require.Equal(t, len(listRes), 2)
		assert.Equal(t, int32(111), listRes[0].ID)
		assert.Equal(t, "title1", listRes[0].Title)
		assert.Equal(t, "desc1", listRes[0].Description)
		assert.Equal(t, int32(0), listRes[0].Position)
		assert.Equal(t, int32(112), listRes[1].ID)
		assert.Equal(t, "title2", listRes[1].Title)
		assert.Equal(t, "desc2", listRes[1].Description)
		assert.Equal(t, int32(1), listRes[1].Position)

		mockedRepo.AssertExpectations(t)
	})
}
