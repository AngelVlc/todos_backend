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

func deleteCategoryRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestDeletesCategoryHandler_Returns_An_Error_If_The_Query_To_Find_The_Existing_Category_Fails(t *testing.T) {
	request := deleteCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	mockedRepo.On("FindCategory", request.Context(), domain.CategoryEntity{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := DeleteCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesCategoryHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Delete_Fails(t *testing.T) {
	request := deleteCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	nvo, _ := domain.NewCategoryNameValueObject("list1")
	existingCategory := domain.CategoryEntity{ID: 11, Name: nvo}
	mockedRepo.On("FindCategory", request.Context(), domain.CategoryEntity{ID: 11, UserID: 1}).Return(&existingCategory, nil).Once()
	mockedRepo.On("DeleteCategory", request.Context(), existingCategory).Return(fmt.Errorf("some error")).Once()

	result := DeleteCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the user category")
	mockedRepo.AssertExpectations(t)
}

func TestDeletesCategoryHandler_Deletes_The_Category(t *testing.T) {
	request := deleteCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	h := handler.Handler{CategoriesRepository: &mockedRepo}

	nvo, _ := domain.NewCategoryNameValueObject("list1")
	existingCategory := domain.CategoryEntity{ID: 11, Name: nvo, UserID: 1}
	mockedRepo.On("FindCategory", request.Context(), domain.CategoryEntity{ID: 11, UserID: 1}).Return(&existingCategory, nil).Once()
	mockedRepo.On("DeleteCategory", request.Context(), existingCategory).Return(nil).Once()

	result := DeleteCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedRepo.AssertExpectations(t)
}
