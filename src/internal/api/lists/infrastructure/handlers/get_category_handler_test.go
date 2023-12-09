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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getCategoryRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestGetCategoryHandler_Returns_An_Error_If_The_Query_To_Find_The_Category_Fails(t *testing.T) {
	request := getCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	mockedRepo.On("FindCategory", request.Context(), domain.CategoryEntity{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := GetCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestGetCategoryHandler_Returns_The_Category(t *testing.T) {
	request := getCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	nvo, _ := domain.NewCategoryNameValueObject("category1")
	foundCategory := domain.CategoryEntity{ID: 11, Name: nvo}
	mockedRepo.On("FindCategory", request.Context(), domain.CategoryEntity{ID: 11, UserID: 1}).Return(&foundCategory, nil).Once()

	result := GetCategoryHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	listRes, isOk := okRes.Content.(*domain.CategoryEntity)
	require.Equal(t, true, isOk, "should be a CategoryEntity")

	assert.Equal(t, int32(11), listRes.ID)
	assert.Equal(t, "category1", listRes.Name.String())
	mockedRepo.AssertExpectations(t)
}
