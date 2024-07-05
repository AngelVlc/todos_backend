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

func getAllCategoriesRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestGetAllCategoriesHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_Fails(t *testing.T) {
	request := getAllCategoriesRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	mockedRepo.On("GetCategories", request.Context(), domain.CategoryRecord{UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := GetAllCategoriesHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting all user categories")
	mockedRepo.AssertExpectations(t)
}

func TestGetAllCategoriesHandler_Returns_The_Categories(t *testing.T) {
	request := getAllCategoriesRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	found := domain.CategoryRecords{
		{ID: 11, Name: "category1"},
		{ID: 12, Name: "category2"},
	}

	mockedRepo.On("GetCategories", request.Context(), domain.CategoryRecord{UserID: 1}).Return(found, nil)

	result := GetAllCategoriesHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	categoriesRes, isOk := okRes.Content.([]*domain.CategoryEntity)
	require.Equal(t, true, isOk, "should be an array of CategoryEntity")

	require.Equal(t, len(categoriesRes), 2)
	assert.Equal(t, int32(11), categoriesRes[0].ID)
	assert.Equal(t, "category1", categoriesRes[0].Name.String())
	assert.Equal(t, int32(12), categoriesRes[1].ID)
	assert.Equal(t, "category2", categoriesRes[1].Name.String())

	mockedRepo.AssertExpectations(t)
}
