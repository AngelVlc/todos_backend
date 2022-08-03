//go:build !e2e
// +build !e2e

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestCreateListItemHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Request_Does_Not_Have_Body(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}
	result := CreateListItemHandler(httptest.NewRecorder(), request(nil), h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateListItemHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Body_Is_Not_A_CreateListItemRequest(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}
	result := CreateListItemHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateListItemHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateListItemRequest_Has_An_Empty_Title(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}
	createReq := createListItemRequest{Title: ""}
	json, _ := json.Marshal(createReq)
	body := bytes.NewBuffer(json)

	result := CreateListItemHandler(httptest.NewRecorder(), request(body), h)

	results.CheckBadRequestErrorResult(t, result, "The item title can not be empty")
}

func TestCreateListItemHandler_Returns_An_Error_If_The_Query_To_Find_The_List_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	request := func() *http.Request {
		createReq := createListItemRequest{Title: "title", Description: "desc"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Returns_An_Error_With_An_UnexpectedError_If_The_Creation_Of_The_ListItem_Fails(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	request := func() *http.Request {
		createReq := createListItemRequest{Title: "title", Description: "desc"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
	mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	listItem := domain.ListItem{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc"}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the list item")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Creates_The_New_ListItem_When_The_List_Does_Not_Have_Any_Items_Yet(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	request := func() *http.Request {
		createReq := createListItemRequest{Title: "title", Description: "desc"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(0)}
	mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	listItem := domain.ListItem{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc", Position: int32(0)}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListItem)
		*arg = domain.ListItem{ID: int32(1)}
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
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	request := func() *http.Request {
		createReq := createListItemRequest{Title: "title", Description: "desc"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(3)}
	mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("GetListItemsMaxPosition", request().Context(), int32(11), int32(1)).Return(int32(-1), fmt.Errorf("some error")).Once()

	result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting the max position")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListItemHandler_Creates_The_New_ListItem_When_The_List_Already_Has_Some_Item(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	mockedEventBus := events.MockedEventBus{}
	h := handler.Handler{ListsRepository: &mockedRepo, EventBus: &mockedEventBus}

	request := func() *http.Request {
		createReq := createListItemRequest{Title: "title", Description: "desc"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(2)}
	mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
	mockedRepo.On("GetListItemsMaxPosition", request().Context(), int32(11), int32(1)).Return(int32(2), nil).Once()
	listItem := domain.ListItem{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc", Position: int32(3)}
	mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.ListItem)
		*arg = domain.ListItem{ID: int32(1)}
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
