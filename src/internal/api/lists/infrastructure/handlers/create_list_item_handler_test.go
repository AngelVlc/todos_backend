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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateListItemHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateListItemInput_Has_An_Empty_Title(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{
		RequestInput: &domain.CreateListItemInput{Title: ""},
	}

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckBadRequestErrorResult(t, result, "The item title can not be empty")
}

func TestCreateListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
		RequestInput:    &domain.CreateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Returns_An_Error_With_An_UnexpectedError_If_The_Creation_Of_The_ListItem_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
		RequestInput:    &domain.CreateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{Name: listName, UserID: int32(1)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	listItem := domain.ListItemRecord{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc"}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the list item")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Creates_The_New_ListItem_When_The_List_Does_Not_Have_Any_Items_Yet(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
		RequestInput:    &domain.CreateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{Name: listName, UserID: int32(1), ItemsCount: int32(0)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	listItem := domain.ListItemRecord{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc", Position: int32(0)}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListItemRecord)
		*arg = domain.ListItemRecord{ID: int32(1)}
	})

	mockedEventBus.On("Publish", "listItemCreated", int32(11))

	mockedEventBus.Wg.Add(1)
	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)
	mockedEventBus.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(infrastructure.ListItemResponse)
	require.Equal(t, true, isOk, "should be a ListItemResponse")
	assert.Equal(t, int32(1), res.ID)

	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}

func TestCreateListItemHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_List_Has_Some_Item_But_GetMaxPosition_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
		RequestInput:    &domain.CreateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{Name: listName, UserID: int32(1), ItemsCount: int32(3)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("GetListItemsMaxPosition", request().Context(), int32(11), int32(1)).Return(int32(-1), fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting the max position")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Creates_The_New_ListItem_When_The_List_Already_Has_Some_Item(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{
		ListsRepository: &mockedRepo,
		EventBus:        &mockedEventBus,
		RequestInput:    &domain.CreateListItemInput{Title: "title", Description: "desc"},
	}

	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListRecord{Name: listName, UserID: int32(1), ItemsCount: int32(2)}
	mockedRepo.On("FindList", request().Context(), &domain.ListRecord{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("GetListItemsMaxPosition", request().Context(), int32(11), int32(1)).Return(int32(2), nil).Once()
	listItem := domain.ListItemRecord{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc", Position: int32(3)}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListItemRecord)
		*arg = domain.ListItemRecord{ID: int32(1)}
	})

	mockedEventBus.On("Publish", "listItemCreated", int32(11))

	mockedEventBus.Wg.Add(1)
	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)
	mockedEventBus.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(infrastructure.ListItemResponse)
	require.Equal(t, true, isOk, "should be a ListItemResponse")
	assert.Equal(t, int32(1), res.ID)

	mockedRepo.AssertExpectations(t)
	mockedEventBus.AssertExpectations(t)
}
