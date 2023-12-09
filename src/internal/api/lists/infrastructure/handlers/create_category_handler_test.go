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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createCategoryRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
	ctx := request.Context()
	ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))

	return request.WithContext(ctx)
}

func TestCreateCategoryHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Check_If_A_Category_With_The_Same_Name_Already_Exists_Fails(t *testing.T) {
	request := createCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	mockedRepo.On("ExistsCategory", request.Context(), domain.CategoryEntity{Name: nvo, UserID: 1}).Return(false, fmt.Errorf("some error")).Once()

	result := CreateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error checking if a category with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestCreateCategoryHandler_Returns_An_Error_Result_With_A_BadRequestError_If_A_Category_With_The_Same_Name_Already_Exists(t *testing.T) {
	request := createCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	mockedRepo.On("ExistsCategory", request.Context(), domain.CategoryEntity{Name: nvo, UserID: 1}).Return(true, nil).Once()

	result := CreateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "A category with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestCreateCategoryHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Creating_The_Category_Fails(t *testing.T) {
	request := createCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	mockedRepo.On("ExistsCategory", request.Context(), domain.CategoryEntity{Name: nvo, UserID: 1}).Return(false, nil).Once()
	newCategory := domain.CategoryEntity{
		Name:   nvo,
		UserID: 1,
	}
	mockedRepo.On("CreateCategory", request.Context(), &newCategory).Return(nil, fmt.Errorf("some error")).Once()

	result := CreateCategoryHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the user category")
	mockedRepo.AssertExpectations(t)
}

func TestCreateCategoryHandler_Creates_A_New_Category(t *testing.T) {
	request := createCategoryRequest()

	mockedRepo := listsRepository.MockedCategoriesRepository{}
	nvo, _ := domain.NewCategoryNameValueObject("category1")
	h := handler.Handler{
		CategoriesRepository: &mockedRepo,
		RequestInput:         &infrastructure.CategoryInput{Name: nvo},
	}

	mockedRepo.On("ExistsCategory", request.Context(), domain.CategoryEntity{Name: nvo, UserID: 1}).Return(false, nil).Once()
	newCategory := domain.CategoryEntity{
		Name:   nvo,
		UserID: 1,
	}
	existingCategory := domain.CategoryEntity{
		ID:     1,
		Name:   nvo,
		UserID: 1,
	}
	mockedRepo.On("CreateCategory", request.Context(), &newCategory).Return(&existingCategory, nil).Once()

	result := CreateCategoryHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(*domain.CategoryEntity)
	require.True(t, isOk, "should be a CategoryEntity")
	assert.Equal(t, int32(1), res.ID)
	assert.Equal(t, "category1", res.Name.String())

	mockedRepo.AssertExpectations(t)
}
