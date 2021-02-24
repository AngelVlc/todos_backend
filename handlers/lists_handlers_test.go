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

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	t.Run("Should return an errorResult if getting the user lists fails", func(t *testing.T) {
		mockedListsService.On("GetUserLists", userID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserListsHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

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

		result := GetUserListsHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.([]*dtos.ListResponseDto)
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

	t.Run("Should return an errorResult if the list id is valid but the query fails", func(t *testing.T) {
		mockedListsService.On("GetUserList", int32(11), userID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

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

		result := GetUserSingleListHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.(*dtos.ListResponseDto)
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

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := AddUserListHandler(httptest.NewRecorder(), request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the insert fails", func(t *testing.T) {
		mockedListsService.On("AddUserList", userID, &dto).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListHandler(httptest.NewRecorder(), request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list if the body is valid and the insert does not fail", func(t *testing.T) {
		mockedListsService.On("AddUserList", userID, &dto).Return(int32(40), nil).Once()

		result := AddUserListHandler(httptest.NewRecorder(), request(true), handler)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(40), okRes.content)

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

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := UpdateUserListHandler(httptest.NewRecorder(), request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		mockedListsService.On("UpdateUserList", int32(40), userID, &dto).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListHandler(httptest.NewRecorder(), request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should update the list if the body is valid and the update does not fail", func(t *testing.T) {
		mockedListsService.On("UpdateUserList", int32(40), userID, &dto).Return(nil).Once()

		result := UpdateUserListHandler(httptest.NewRecorder(), request(true), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.(*dtos.ListDto)
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

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		mockedListsService.On("RemoveUserList", int32(40), userID).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserListHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list if there is no errors", func(t *testing.T) {
		mockedListsService.On("RemoveUserList", int32(40), userID).Return(nil).Once()

		result := DeleteUserListHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, okRes.content)

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

	t.Run("Should return an errorResult if the listId is valid but the query fails", func(t *testing.T) {
		res := dtos.GetItemResultDto{}
		mockedListsService.On("GetUserListItem", int32(3), int32(5), userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListItemHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return the item it there is no errors", func(t *testing.T) {
		res := dtos.GetItemResultDto{
			Title:       "Title",
			Description: "Description",
		}
		mockedListsService.On("GetUserListItem", int32(3), int32(5), userID, &dtos.GetItemResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(3).(*dtos.GetItemResultDto)
			*arg = res
		})

		result := GetUserSingleListItemHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.(dtos.GetItemResultDto)
		require.Equal(t, true, isOk, "should be an item result dto")
		assert.Equal(t, res.Title, resDto.Title)
		assert.Equal(t, res.Description, resDto.Description)

		mockedListsService.AssertExpectations(t)
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

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := AddUserListItemHandler(httptest.NewRecorder(), request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the insert fails", func(t *testing.T) {
		mockedListsService.On("AddUserListItem", listID, userID, &dto).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListItemHandler(httptest.NewRecorder(), request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list item if the body is valid and the insert does not fail", func(t *testing.T) {
		mockedListsService.On("AddUserListItem", listID, userID, &dto).Return(int32(40), nil).Once()

		result := AddUserListItemHandler(httptest.NewRecorder(), request(true), handler)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(40), okRes.content)

		mockedListsService.AssertExpectations(t)
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

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		mockedListsService.On("RemoveUserListItem", int32(20), int32(40), userID).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserListItemHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list item if there is no errors", func(t *testing.T) {
		mockedListsService.On("RemoveUserListItem", int32(20), int32(40), userID).Return(nil).Once()

		result := DeleteUserListItemHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, okRes.content)

		mockedListsService.AssertExpectations(t)
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

	t.Run("Should return an errorResult if the body is not valid", func(t *testing.T) {
		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		mockedListsService.On("UpdateUserListItem", listID, int32(40), userID, &dto).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should update the list item if the body is valid and the update does not fail", func(t *testing.T) {
		mockedListsService.On("UpdateUserListItem", listID, int32(40), userID, &dto).Return(nil).Once()

		result := UpdateUserListItemHandler(httptest.NewRecorder(), request(true), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.(*dtos.ListItemDto)
		require.Equal(t, true, isOk, "should be a list item result dto")
		assert.Equal(t, dto.Title, resDto.Title)

		mockedListsService.AssertExpectations(t)
	})
}
