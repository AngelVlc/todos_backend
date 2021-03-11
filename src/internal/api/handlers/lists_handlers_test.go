//+build !e2e

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/services"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userID = int32(21)
)

func TestGetUserListsHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}
	mockedListsService := services.NewMockedListsService()
	h := handler.Handler{ListsSrv: mockedListsService}

	t.Run("Should return an errorResult if getting the user lists fails", func(t *testing.T) {
		mockedListsService.On("GetUserLists", userID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserListsHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the user lists if there is no errors", func(t *testing.T) {
		res := []*dtos.ListResponseDto{
			{
				ID:   int32(1),
				Name: "list1",
			},
		}
		mockedListsService.On("GetUserLists", userID).Return(res, nil).Once()

		result := GetUserListsHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.([]*dtos.ListResponseDto)
		require.Equal(t, true, isOk, "should be a list response dto")
		require.Equal(t, 1, len(resDto))
		assert.Equal(t, res[0].ID, resDto[0].ID)
		assert.Equal(t, res[0].Name, resDto[0].Name)

		mockedListsService.AssertExpectations(t)
	})
}

func TestGetUserSingleListHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListsService := services.NewMockedListsService()
	h := handler.Handler{ListsSrv: mockedListsService}

	t.Run("Should return an errorResult if the list id is valid but the query fails", func(t *testing.T) {
		mockedListsService.On("GetUserList", int32(11), userID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return the list it there is no errors", func(t *testing.T) {
		res := dtos.ListResponseDto{
			ID:   int32(11),
			Name: "list1",
			ListItems: []*dtos.ListItemResponseDto{
				{
					ID:          int32(1),
					Title:       "the title",
					Description: "the description",
				},
			},
		}
		mockedListsService.On("GetUserList", int32(11), userID).Return(&res, nil).Once()

		result := GetUserSingleListHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*dtos.ListResponseDto)
		require.Equal(t, true, isOk, "should be a list response dto")
		assert.Equal(t, res.ID, resDto.ID)
		assert.Equal(t, res.Name, resDto.Name)
		require.Equal(t, 1, len(resDto.ListItems))
		assert.Equal(t, res.ListItems[0].ID, resDto.ListItems[0].ID)
		assert.Equal(t, res.ListItems[0].Title, resDto.ListItems[0].Title)
		assert.Equal(t, res.ListItems[0].Description, resDto.ListItems[0].Description)

		mockedListsService.AssertExpectations(t)
	})
}

func TestAddUserListHandler(t *testing.T) {
	dto := dtos.ListDto{
		Name: "list",
	}
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			json, _ := json.Marshal(dto)
			body = bytes.NewBuffer(json)
		} else {
			body = strings.NewReader("wadus")
		}
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListsService := services.NewMockedListsService()
	h := handler.Handler{ListsSrv: mockedListsService}

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := AddUserListHandler(httptest.NewRecorder(), request(false), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the create fails", func(t *testing.T) {
		mockedListsService.On("AddUserList", userID, &dto).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListHandler(httptest.NewRecorder(), request(true), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list if the body is valid and the create does not fail", func(t *testing.T) {
		mockedListsService.On("AddUserList", userID, &dto).Return(int32(40), nil).Once()

		result := AddUserListHandler(httptest.NewRecorder(), request(true), h)

		okRes := results.CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(40), okRes.Content)

		mockedListsService.AssertExpectations(t)
	})
}

func TestUpdateUserListHandler(t *testing.T) {
	dto := dtos.ListDto{
		Name: "list",
	}
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			json, _ := json.Marshal(dto)
			body = bytes.NewBuffer(json)
		} else {
			body = strings.NewReader("wadus")
		}
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListsService := services.NewMockedListsService()
	h := handler.Handler{ListsSrv: mockedListsService}

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := UpdateUserListHandler(httptest.NewRecorder(), request(false), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		mockedListsService.On("UpdateUserList", int32(40), userID, &dto).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListHandler(httptest.NewRecorder(), request(true), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should update the list if the body is valid and the update does not fail", func(t *testing.T) {
		mockedListsService.On("UpdateUserList", int32(40), userID, &dto).Return(nil).Once()

		result := UpdateUserListHandler(httptest.NewRecorder(), request(true), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*dtos.ListDto)
		require.Equal(t, true, isOk, "should be a list")
		assert.Equal(t, resDto.Name, dto.Name)

		mockedListsService.AssertExpectations(t)
	})
}

func TestDeleteUserListHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListsService := services.NewMockedListsService()
	h := handler.Handler{ListsSrv: mockedListsService}

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		mockedListsService.On("RemoveUserList", int32(40), userID).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserListHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list if there is no errors", func(t *testing.T) {
		mockedListsService.On("RemoveUserList", int32(40), userID).Return(nil).Once()

		result := DeleteUserListHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, okRes.Content)

		mockedListsService.AssertExpectations(t)
	})
}

func TestGetUserSingleListItemHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "5",
			"itemId": "3",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListItemsService := services.NewMockedListItemsService()
	h := handler.Handler{ListItemsSrv: mockedListItemsService}

	t.Run("Should return an errorResult if the listId is valid but the query fails", func(t *testing.T) {
		mockedListItemsService.On("GetListItem", int32(3), int32(5), userID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListItemsService.AssertExpectations(t)
	})

	t.Run("Should return the item it there is no errors", func(t *testing.T) {
		res := dtos.ListItemResponseDto{
			Title:       "Title",
			Description: "Description",
		}
		mockedListItemsService.On("GetListItem", int32(3), int32(5), userID).Return(&res, nil).Once()

		result := GetUserSingleListItemHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*dtos.ListItemResponseDto)
		require.Equal(t, true, isOk, "should be a list item response dto")
		assert.Equal(t, res.Title, resDto.Title)
		assert.Equal(t, res.Description, resDto.Description)

		mockedListItemsService.AssertExpectations(t)
	})
}

func TestAddUserListItemHandler(t *testing.T) {
	listID := int32(11)
	dto := dtos.ListItemDto{
		Title: "title",
	}
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			json, _ := json.Marshal(dto)
			body = bytes.NewBuffer(json)
		} else {
			body = strings.NewReader("wadus")
		}
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"listId": "11",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListsService := services.NewMockedListsService()
	mockedListItemsService := services.NewMockedListItemsService()
	h := handler.Handler{ListsSrv: mockedListsService, ListItemsSrv: mockedListItemsService}

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := AddUserListItemHandler(httptest.NewRecorder(), request(false), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the create fails", func(t *testing.T) {
		mockedListItemsService.On("AddListItem", listID, userID, &dto).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListItemHandler(httptest.NewRecorder(), request(true), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListItemsService.AssertExpectations(t)
	})

	t.Run("Should add the list item if the body is valid and the create does not fail", func(t *testing.T) {
		mockedListItemsService.On("AddListItem", listID, userID, &dto).Return(int32(40), nil).Once()

		result := AddUserListItemHandler(httptest.NewRecorder(), request(true), h)

		okRes := results.CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(40), okRes.Content)

		mockedListItemsService.AssertExpectations(t)
	})
}

func TestDeleteUserListItemHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"itemId": "20",
			"listId": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListItemsService := services.NewMockedListItemsService()
	h := handler.Handler{ListItemsSrv: mockedListItemsService}

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		mockedListItemsService.On("RemoveListItem", int32(20), int32(40), userID).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserListItemHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListItemsService.AssertExpectations(t)
	})

	t.Run("Should delete the list item if there is no errors", func(t *testing.T) {
		mockedListItemsService.On("RemoveListItem", int32(20), int32(40), userID).Return(nil).Once()

		result := DeleteUserListItemHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, okRes.Content)

		mockedListItemsService.AssertExpectations(t)
	})
}

func TestUpdateUserListItemHandler(t *testing.T) {
	listID := int32(20)
	dto := dtos.ListItemDto{
		Title: "title",
	}
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			json, _ := json.Marshal(dto)
			body = bytes.NewBuffer(json)
		} else {
			body = strings.NewReader("wadus")
		}
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"itemId": "20",
			"listId": "40",
		})
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, userID)
		return request.WithContext(ctx)
	}

	mockedListItemsService := services.NewMockedListItemsService()
	h := handler.Handler{ListItemsSrv: mockedListItemsService}

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(false), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		mockedListItemsService.On("UpdateListItem", listID, int32(40), userID, &dto).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(true), h)

		results.CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListItemsService.AssertExpectations(t)
	})

	t.Run("Should update the list item if the body is valid and the update does not fail", func(t *testing.T) {
		mockedListItemsService.On("UpdateListItem", listID, int32(40), userID, &dto).Return(nil).Once()

		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(true), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*dtos.ListItemDto)
		require.Equal(t, true, isOk, "should be a list item result dto")
		assert.Equal(t, dto.Title, resDto.Title)

		mockedListItemsService.AssertExpectations(t)
	})
}
