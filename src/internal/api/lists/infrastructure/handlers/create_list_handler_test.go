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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateListHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Request_Does_Not_Have_Body(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}

	result := CreateListHandler(httptest.NewRecorder(), request(nil), h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateListHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Body_Is_Not_A_CreateListRequest(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}

	result := CreateListHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

	results.CheckBadRequestErrorResult(t, result, "Invalid body")
}

func TestCreateListHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_The_CreateListRequest_Has_An_Empty_Name(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, int32(1))
		return request.WithContext(ctx)
	}

	h := handler.Handler{}
	createReq := createListRequest{Name: ""}
	json, _ := json.Marshal(createReq)
	body := bytes.NewBuffer(json)

	result := CreateListHandler(httptest.NewRecorder(), request(body), h)

	results.CheckBadRequestErrorResult(t, result, "The list name can not be empty")
}

func TestCreateListHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Query_To_Check_If_A_List_With_The_Same_Name_Already_Exists(t *testing.T) {
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

	mockedRepo.On("ExistsList", request().Context(), &domain.List{Name: domain.ListName("list1"), UserID: int32(1)}).Return(false, fmt.Errorf("some error")).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Returns_An_Error_Result_With_A_BadRequestError_If_A_List_With_The_Same_Name_Already_Exists(t *testing.T) {
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

	mockedRepo.On("ExistsList", request().Context(), &domain.List{Name: domain.ListName("list1"), UserID: int32(1)}).Return(true, nil).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckBadRequestErrorResult(t, result, "A list with the same name already exists")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Returns_An_Error_Result_With_An_UnexpectedError_If_Creating_The_List_Fails(t *testing.T) {
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

	mockedRepo.On("ExistsList", request().Context(), &domain.List{Name: domain.ListName("list1"), UserID: int32(1)}).Return(false, nil).Once()
	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
	mockedRepo.On("CreateList", request().Context(), &list).Return(fmt.Errorf("some error")).Once()

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating the user list")
	mockedRepo.AssertExpectations(t)
}

func TestCreateListHandler_Creates_A_New_List(t *testing.T) {
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

	mockedRepo.On("ExistsList", request().Context(), &domain.List{Name: domain.ListName("list1"), UserID: int32(1)}).Return(false, nil).Once()
	list := domain.List{Name: domain.ListName("list1"), UserID: int32(1)}
	mockedRepo.On("CreateList", request().Context(), &list).Return(nil).Once().Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.List)
		*arg = domain.List{ID: int32(1), Name: "list1"}
	})

	result := CreateListHandler(httptest.NewRecorder(), request(), h)

	okRes := results.CheckOkResult(t, result, http.StatusCreated)
	res, isOk := okRes.Content.(infrastructure.ListResponse)
	require.True(t, isOk, "should be a ListResponse")
	assert.Equal(t, int32(1), res.ID)
	assert.Equal(t, "list1", res.Name)

	mockedRepo.AssertExpectations(t)
}
