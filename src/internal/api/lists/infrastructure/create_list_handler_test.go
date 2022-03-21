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
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateListHandlerValidations(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		result := CreateListHandler(httptest.NewRecorder(), request(nil), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a create list request", func(t *testing.T) {
		result := CreateListHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create list request has an empty Name", func(t *testing.T) {
		createReq := createListRequest{Name: ""}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		result := CreateListHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The list name can not be empty")
	})
}

func TestCreateListHandler(t *testing.T) {
	mockedRepo := listsRepository.MockedListsRepository{}
	h := handler.Handler{ListsRepository: &mockedRepo}

	request := func() *http.Request {
		createReq := createListRequest{Name: "list1"}
		json, _ := json.Marshal(createReq)
		body := bytes.NewBuffer(json)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	t.Run("Should return an error result with an UnexpectedError if the query to check if a user with the same name exists fails", func(t *testing.T) {
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, fmt.Errorf("some error")).Once()

		result := CreateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with a BadRequestError if a list with the same name already exist", func(t *testing.T) {
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(true, nil).Once()

		result := CreateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "A list with the same name already exists")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if create user list fails", func(t *testing.T) {
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, nil).Once()
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
		mockedRepo.On("CreateList", request().Context(), &list).Return(fmt.Errorf("some error")).Once()

		result := CreateListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating the user list")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should create the new user list", func(t *testing.T) {
		mockedRepo.On("ExistsList", request().Context(), domain.ListName("list1"), int32(1)).Return(false, nil).Once()
		list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
		mockedRepo.On("CreateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*domain.List)
			*arg = domain.List{ID: int32(1)}
		})

		result := CreateListHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusCreated)
		res, isOk := okRes.Content.(ListResponse)
		require.Equal(t, true, isOk, "should be a ListResponse")
		assert.Equal(t, int32(1), res.ID)

		mockedRepo.AssertExpectations(t)
	})
}
