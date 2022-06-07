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
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateListHandlerValidations(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		result := UpdateListHandler(httptest.NewRecorder(), request(nil), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a update list request", func(t *testing.T) {
		result := UpdateListHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the update list request has an empty Name", func(t *testing.T) {
		updateReq := updateListRequest{Name: ""}
		json, _ := json.Marshal(updateReq)
		body := bytes.NewBuffer(json)

		result := UpdateListHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The list name can not be empty")
	})
}

func TestUpdateListHandler(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	request := func() *http.Request {
		updateReq := updateListRequest{Name: "list1", IDsByPosition: []int32{int32(2), int32(1)}}
		json, _ := json.Marshal(updateReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	t.Run("Should return an error  if the query to find the list fails", func(t *testing.T) {
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with an UnexpectedError if it is trying to update the list name but the query to find a list by name fails", func(t *testing.T) {
		list := domain.List{ID: int32(11), Name: domain.ListName("oldName"), UserID: int32(11)}
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, fmt.Errorf("some error")).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if updating the user list fails", func(t *testing.T) {
		list := domain.List{ID: int32(11), Name: domain.ListName("list1"), UserID: int32(11)}
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
		mockedRepo.On("UpdateList", request().Context(), &list).Return(fmt.Errorf("some error")).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error updating the user list")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if getting the user list items fails", func(t *testing.T) {
		list := domain.List{ID: int32(11), Name: domain.ListName("originalName"), UserID: int32(1)}
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, nil).Once()
		mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.List)
			*arg = domain.List{Name: "list1"}
		})
		mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting all list items")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if bulk updating the list items fails", func(t *testing.T) {
		list := domain.List{ID: int32(11), Name: domain.ListName("originalName"), UserID: int32(1)}
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, nil).Once()
		mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.List)
			*arg = domain.List{Name: "list1"}
		})
		listItems := []domain.ListItem{{ID: int32(1), Position: int32(0)}, {ID: int32(2), Position: int32(1)}}
		mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(listItems, nil).Once()
		mockedRepo.On("BulkUpdateListItems", request().Context(), listItems).Return(fmt.Errorf("some error")).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error bulk updating")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should update the user list and the items position", func(t *testing.T) {
		list := domain.List{ID: int32(11), Name: domain.ListName("originalName"), UserID: int32(1)}
		mockedRepo.On("FindList", request().Context(), &domain.List{ID: int32(11), UserID: int32(1)}).Return(&list, nil).Once()
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, nil).Once()
		mockedRepo.On("UpdateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.List)
			*arg = domain.List{Name: "list1"}
		})
		listItems := []domain.ListItem{{ID: int32(1), Position: int32(0)}, {ID: int32(2), Position: int32(1)}}
		mockedRepo.On("GetAllListItems", request().Context(), list.ID, list.UserID).Return(listItems, nil).Once()
		listItems = []domain.ListItem{{ID: int32(1), Position: int32(1)}, {ID: int32(2), Position: int32(0)}}
		mockedRepo.On("BulkUpdateListItems", request().Context(), listItems).Return(nil).Once()

		result := UpdateListHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		res, isOk := okRes.Content.(infrastructure.ListResponse)
		require.Equal(t, true, isOk, "should be a ListResponse")
		assert.Equal(t, "list1", res.Name)

		mockedRepo.AssertExpectations(t)
	})
}
