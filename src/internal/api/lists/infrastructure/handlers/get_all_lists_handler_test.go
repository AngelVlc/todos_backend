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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAllRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestGetAllListsHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_Fails(t *testing.T) {
	request := getAllRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	mockedRepo.On("GetAllListsForUser", request.Context(), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

	result := GetAllListsHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting all user lists")
	mockedRepo.AssertExpectations(t)
}

func TestGetAllListsHandler_Returns_The_Lists(t *testing.T) {
	request := getAllRequest()

	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	found := []domain.ListRecord{
		{ID: 11, Name: "list1", ItemsCount: 4},
		{ID: 12, Name: "list2", ItemsCount: 8},
	}

	mockedRepo.On("GetAllListsForUser", request.Context(), int32(1)).Return(found, nil)

	result := GetAllListsHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	listRes, isOk := okRes.Content.([]*domain.ListEntity)
	require.Equal(t, true, isOk, "should be an array of ListEntity")

	require.Equal(t, len(listRes), 2)
	assert.Equal(t, int32(11), listRes[0].ID)
	assert.Equal(t, "list1", listRes[0].Name.String())
	assert.Equal(t, int32(4), listRes[0].ItemsCount)
	assert.Equal(t, int32(12), listRes[1].ID)
	assert.Equal(t, "list2", listRes[1].Name.String())
	assert.Equal(t, int32(8), listRes[1].ItemsCount)

	mockedRepo.AssertExpectations(t)
}
