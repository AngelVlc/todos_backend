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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllListsHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	t.Run("Should return an error result with an unexpected error if the query fails", func(t *testing.T) {
		mockedRepo.On("GetAllLists", int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetAllListsHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting all user lists")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return the lists if the query does not fail", func(t *testing.T) {
		found := []domain.List{
			{ID: 11, Name: "list1", ItemsCount: 4},
			{ID: 12, Name: "list2", ItemsCount: 8},
		}

		mockedRepo.On("GetAllLists", int32(1)).Return(found, nil)

		result := GetAllListsHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		listRes, isOk := okRes.Content.([]ListResponse)
		require.Equal(t, true, isOk, "should be an array of list response")

		require.Equal(t, len(listRes), 2)
		assert.Equal(t, int32(11), listRes[0].ID)
		assert.Equal(t, "list1", listRes[0].Name)
		assert.Equal(t, int32(4), listRes[0].ItemsCount)
		assert.Equal(t, int32(12), listRes[1].ID)
		assert.Equal(t, "list2", listRes[1].Name)
		assert.Equal(t, int32(8), listRes[1].ItemsCount)

		mockedRepo.AssertExpectations(t)
	})
}
