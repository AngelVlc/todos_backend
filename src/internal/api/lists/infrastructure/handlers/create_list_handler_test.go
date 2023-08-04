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
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Check_If_A_List_With_The_Same_Name_Already_Exists_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	mockedRepo.On("ExistsList", request().Context(), domain.ListEntity{Name: listName, UserID: 1}).Return(false, fmt.Errorf("some error")).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error checking if a list with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Returns_An_Error_Result_With_A_BadRequestError_If_A_List_With_The_Same_Name_Already_Exists(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	mockedRepo.On("ExistsList", request().Context(), domain.ListEntity{Name: listName, UserID: 1}).Return(true, nil).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckBadRequestErrorResult(t, result, "A list with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Creating_The_List_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
	}

	mockedRepo.On("ExistsList", request().Context(), domain.ListEntity{Name: listName, UserID: 1}).Return(false, nil).Once()
	createdList := domain.ListEntity{
		Name:   listName,
		UserID: 1,
		Items:  []*domain.ListItemEntity{},
	}
	mockedRepo.On("CreateList", request().Context(), &createdList).Return(nil, fmt.Errorf("some error")).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the user list")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Creates_A_New_List_And_Sends_The_ListCreatedOrUpdated_Event(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	listName, _ := domain.NewListNameValueObject("list1")
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		RequestInput:    &infrastructure.ListInput{Name: listName},
		EventBus:        &mockedEventBus,
	}

	mockedRepo.On("ExistsList", request().Context(), domain.ListEntity{Name: listName, UserID: 1}).Return(false, nil).Once()
	listToCreate := domain.ListEntity{
		Name:   listName,
		UserID: 1,
		Items:  []*domain.ListItemEntity{},
	}
	createdList := domain.ListEntity{
		ID:     1,
		Name:   listName,
		UserID: 1,
		Items:  []*domain.ListItemEntity{},
	}
	mockedRepo.On("CreateList", request().Context(), &listToCreate).Return(&createdList, nil).Once()

	mockedEventBus.On("Publish", "listCreatedOrUpdated", int32(1))

	mockedEventBus.Wg.Add(1)
	result := CreateListHandler(httptest.NewRecorder(), request(), h)
	mockedEventBus.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(*domain.ListEntity)
	require.True(t, isOk, "should be a ListEntity")
	assert.Equal(t, int32(1), res.ID)
	assert.Equal(t, "list1", res.Name.String())

	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
