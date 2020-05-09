package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserListsHandler(t *testing.T) {
	mockedListsService := services.NewMockedListsService()

	handler := Handler{
		listsSrv: mockedListsService,
	}

	userID := int32(21)

	t.Run("Should return an errorResult if getting the user lists fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := []dtos.GetListsResultDto{}
		mockedListsService.On("GetUserLists", userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserListsHandler(request.WithContext(ctx), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the user lists if there is no errors", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := []dtos.GetListsResultDto{
			dtos.GetListsResultDto{
				ID:   int32(1),
				Name: "list1",
			},
		}
		mockedListsService.On("GetUserLists", userID, &[]dtos.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]dtos.GetListsResultDto)
			*arg = res
		})

		result := GetUserListsHandler(request.WithContext(ctx), handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestGetUserSingleListHandler(t *testing.T) {
	mockedListsService := services.NewMockedListsService()

	handler := Handler{
		listsSrv: mockedListsService,
	}

	userID := int32(21)

	t.Run("Should return an errorResult if user id url param is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "badId",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		result := GetUserSingleListHandler(request.WithContext(ctx), handler)

		CheckBadRequestErrorResult(t, result, "Invalid id in url")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the id is valid but the query fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := dtos.GetSingleListResultDto{}
		mockedListsService.On("GetSingleUserList", int32(11), userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListHandler(request.WithContext(ctx), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return the list it there is no errors", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := dtos.GetSingleListResultDto{
			ID:   int32(11),
			Name: "list1",
			ListItems: []dtos.GetSingleListResultItemDto{
				dtos.GetSingleListResultItemDto{
					ID:          int32(1),
					ListID:      int32(11),
					Title:       "the title",
					Description: "the description",
				},
			},
		}
		mockedListsService.On("GetSingleUserList", int32(11), userID, &dtos.GetSingleListResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(2).(*dtos.GetSingleListResultDto)
			*arg = res
		})

		result := GetUserSingleListHandler(request.WithContext(ctx), handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestAddUserListHandler(t *testing.T) {
	mockedListsService := services.NewMockedListsService()

	handler := Handler{
		listsSrv: mockedListsService,
	}

	userID := int32(21)

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		result := AddUserListHandler(request.WithContext(ctx), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the insert fails", func(t *testing.T) {
		dto := dtos.ListDto{
			Name: "list",
		}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := models.List{Name: "list"}
		mockedListsService.On("AddUserList", userID, &res).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListHandler(request.WithContext(ctx), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list if the body is valid and the insert does not fail", func(t *testing.T) {
		dto := dtos.ListDto{
			Name: "list",
		}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := models.List{Name: "list"}
		mockedListsService.On("AddUserList", userID, &res).Return(int32(40), nil).Once()

		result := AddUserListHandler(request.WithContext(ctx), handler)

		assert.Equal(t, okResult{int32(40), http.StatusCreated}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestUpdateUserListHandler(t *testing.T) {
	mockedListsService := services.NewMockedListsService()

	handler := Handler{
		listsSrv: mockedListsService,
	}

	userID := int32(21)

	t.Run("Should return an errorResult if user id url param is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "badId",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		result := UpdateUserListHandler(request.WithContext(ctx), handler)

		CheckBadRequestErrorResult(t, result, "Invalid id in url")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		result := UpdateUserListHandler(request.WithContext(ctx), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		dto := dtos.ListDto{
			Name: "list",
		}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := models.List{Name: "list"}
		mockedListsService.On("UpdateUserList", int32(40), userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListHandler(request.WithContext(ctx), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should update the list if the body is valid and the update does not fail", func(t *testing.T) {
		dto := dtos.ListDto{
			Name: "list",
		}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		res := models.List{Name: "list"}
		mockedListsService.On("UpdateUserList", int32(40), userID, &res).Return(nil).Once()

		result := UpdateUserListHandler(request.WithContext(ctx), handler)

		assert.Equal(t, okResult{&res, http.StatusOK}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestDeleteUserListHandler(t *testing.T) {
	mockedListsService := services.NewMockedListsService()

	handler := Handler{
		listsSrv: mockedListsService,
	}

	userID := int32(21)

	t.Run("Should return an errorResult if user id url param is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "badId",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		result := DeleteUserListHandler(request.WithContext(ctx), handler)

		CheckBadRequestErrorResult(t, result, "Invalid id in url")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the update fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		mockedListsService.On("RemoveUserList", int32(40), userID).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserListHandler(request.WithContext(ctx), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list if there is no errors", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)

		mockedListsService.On("RemoveUserList", int32(40), userID).Return(nil).Once()

		result := DeleteUserListHandler(request.WithContext(ctx), handler)

		assert.Equal(t, okResult{nil, http.StatusNoContent}, result)

		mockedListsService.AssertExpectations(t)
	})
}
