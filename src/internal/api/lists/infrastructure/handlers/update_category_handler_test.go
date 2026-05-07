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
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateCategoryRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	request = mux.SetURLVars(request, map[string]string{
		"id": "11",
	})
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestUpdateCategoryHandler_Returns_An_Error_If_The_Query_To_Find_The_Category_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	request := updateCategoryRequest()

	mockedRepo.On("FindCategory", request.Context(), domain.CategoryRecord{ID: 11, UserID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := UpdateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateCategoryHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Is_Trying_To_Update_The_Category_Name_But_The_Query_To_Check_If_The_Name_Exists_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	request := updateCategoryRequest()

	mockedRepo.On("FindCategory", request.Context(), domain.CategoryRecord{ID: 11, UserID: 1}).Return(&domain.CategoryRecord{ID: 11, Name: "oldName"}, nil).Once()
	mockedRepo.On("ExistsCategory", request.Context(), domain.CategoryRecord{Name: "category1", UserID: 1}).Return(false, fmt.Errorf("some error")).Once()

	result := UpdateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "Error checking if a category with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateCategoryHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Updating_The_Category_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	request := updateCategoryRequest()

	category := domain.CategoryRecord{
		ID:     int32(11),
		Name:   "category1",
		UserID: 1,
	}
	mockedRepo.On("FindCategory", request.Context(), domain.CategoryRecord{ID: 11, UserID: 1}).Return(&domain.CategoryRecord{ID: 11, Name: "category1"}, nil).Once()
	mockedRepo.On("UpdateCategory", request.Context(), &category).Return(fmt.Errorf("some error")).Once()

	result := UpdateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error updating the user category")
	mockedRepo.AssertExpectations(t)
}

func TestUpdateCategoryHandler_Updates_The_Category(t *testing.T) {
	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo, IsFavourite: false},
	}

	request := updateCategoryRequest()

	recordToUpdate := domain.CategoryRecord{
		ID:          11,
		Name:        "category1",
		UserID:      1,
		IsFavourite: false,
	}
	mockedRepo.On("FindCategory", request.Context(), domain.CategoryRecord{ID: 11, UserID: 1}).Return(&recordToUpdate, nil).Once()
	mockedRepo.On("UpdateCategory", request.Context(), &recordToUpdate).Return(nil).Once()

	result := UpdateCategoryHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(*domain.CategoryEntity)
	require.True(t, isOk, "should be a CategoryEntity")
	assert.Equal(t, "category1", res.Name.String())

	mockedRepo.AssertExpectations(t)
}
