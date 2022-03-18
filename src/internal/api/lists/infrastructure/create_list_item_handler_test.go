//go:build !e2e
// +build !e2e

package infrastructure

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

	"github.com/AngelVlc/todos_backend/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateListItemHandlerValidations(t *testing.T) {
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

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		result := CreateListItemHandler(httptest.NewRecorder(), request(nil), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a create list item request", func(t *testing.T) {
		result := CreateListItemHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create list item request has an empty Title", func(t *testing.T) {
		createReq := createListItemRequest{Title: ""}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		result := CreateListItemHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The item title can not be empty")
	})
}

func TestCreateListItemHandler(t *testing.T) {
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

	t.Run("Should return an error if the query to find the list fails", func(t *testing.T) {
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if create list item fails", func(t *testing.T) {
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(&list, nil).Once()
		listItem := domain.ListItem{ListID: int32(11), UserID: int32(1), Title: "title", Description: "desc"}
		mockedRepo.On("CreateListItem", request().Context(), &listItem).Return(fmt.Errorf("some error")).Once()

		result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating the list item")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should create the new list item when the list does not have any items yet", func(t *testing.T) {
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(0)}
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(&list, nil).Once()
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
		res, isOk := okRes.Content.(ListItemResponse)
		require.Equal(t, true, isOk, "should be a ListItemResponse")
		assert.Equal(t, int32(1), res.ID)

		mockedRepo.AssertExpectations(t)
		mockedEventBus.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if the list has some item but get max position fails", func(t *testing.T) {
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(3)}
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(&list, nil).Once()
		mockedRepo.On("GetListItemsMaxPosition", request().Context(), int32(11), int32(1)).Return(int32(-1), fmt.Errorf("some error")).Once()

		result := CreateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting the max position")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should create the new list item when the list already has some item", func(t *testing.T) {
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1), ItemsCount: int32(2)}
		mockedRepo.On("FindListByID", request().Context(), int32(11), int32(1)).Return(&list, nil).Once()
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
		res, isOk := okRes.Content.(ListItemResponse)
		require.Equal(t, true, isOk, "should be a ListItemResponse")
		assert.Equal(t, int32(1), res.ID)

		mockedRepo.AssertExpectations(t)
		mockedEventBus.AssertExpectations(t)
	})
}
