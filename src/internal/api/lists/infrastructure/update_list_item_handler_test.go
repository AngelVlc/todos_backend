//+build !e2e

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

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateListItemHandlerValidations(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		result := UpdateListItemHandler(httptest.NewRecorder(), request(nil), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a update list item request", func(t *testing.T) {
		result := UpdateListItemHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the update list item request has an empty Title", func(t *testing.T) {
		updateReq := updateListItemRequest{Title: ""}
		json, _ := json.Marshal(updateReq)
		body := bytes.NewBuffer(json)

		result := UpdateListItemHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The item title can not be empty")
	})
}

func TestUpdateListItemHandler(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	request := func() *http.Request {
		updateReq := updateListItemRequest{Title: "title"}
		json, _ := json.Marshal(updateReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"id":     "111",
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the list item fails", func(t *testing.T) {
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting the list item")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError if the list item does not exits", func(t *testing.T) {
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(nil, nil).Once()

		result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "The list item does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if update list item fails", func(t *testing.T) {
		listItem := domain.ListItem{ID: 1, ListID: 11, Title: domain.ItemTitle("title"), Description: "desc", UserID: int32(1)}
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(&listItem, nil).Once()
		mockedRepo.On("UpdateListItem", &listItem).Return(fmt.Errorf("some error")).Once()

		result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error updating the list item")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should update the list item", func(t *testing.T) {
		listItem := domain.ListItem{ID: 1, ListID: 11, Title: domain.ItemTitle("oriTitle"), Description: "oriDesc", UserID: int32(1)}
		mockedRepo.On("FindListItemByID", int32(111), int32(11), int32(1)).Return(&listItem, nil).Once()
		mockedRepo.On("UpdateListItem", &listItem).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.ListItem)
			*arg = domain.ListItem{Title: "title", Description: "desc"}
		})

		result := UpdateListItemHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		res, isOk := okRes.Content.(ListItemResponse)
		require.Equal(t, true, isOk, "should be a ListItemResponse")
		assert.Equal(t, "title", res.Title)
		assert.Equal(t, "desc", res.Description)

		mockedRepo.AssertExpectations(t)
	})
}
