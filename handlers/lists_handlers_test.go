package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		res := []dtos.GetListsResultDto{}
		mockedListsService.On("GetUserLists", userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserListsHandler(request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the user lists if there is no errors", func(t *testing.T) {
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

		result := GetUserListsHandler(request(), handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

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
		res := dtos.GetSingleListResultDto{}
		mockedListsService.On("GetSingleUserList", int32(11), userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUserSingleListHandler(request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return the list it there is no errors", func(t *testing.T) {
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

		result := GetUserSingleListHandler(request(), handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestAddUserListHandler(t *testing.T) {
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			dto := dtos.ListDto{
				Name: "list",
			}
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
		result := AddUserListHandler(request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the insert fails", func(t *testing.T) {
		res := models.List{Name: "list"}
		mockedListsService.On("AddUserList", userID, &res).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListHandler(request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list if the body is valid and the insert does not fail", func(t *testing.T) {
		res := models.List{Name: "list"}
		mockedListsService.On("AddUserList", userID, &res).Return(int32(40), nil).Once()

		result := AddUserListHandler(request(true), handler)

		assert.Equal(t, okResult{int32(40), http.StatusCreated}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestUpdateUserListHandler(t *testing.T) {
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			dto := dtos.ListDto{
				Name: "list",
			}
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
		result := UpdateUserListHandler(request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the update fails", func(t *testing.T) {
		res := models.List{Name: "list"}
		mockedListsService.On("UpdateUserList", int32(40), userID, &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := UpdateUserListHandler(request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should update the list if the body is valid and the update does not fail", func(t *testing.T) {
		res := models.List{Name: "list"}
		mockedListsService.On("UpdateUserList", int32(40), userID, &res).Return(nil).Once()

		result := UpdateUserListHandler(request(true), handler)

		assert.Equal(t, okResult{&res, http.StatusOK}, result)

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

		result := DeleteUserListHandler(request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list if there is no errors", func(t *testing.T) {
		mockedListsService.On("RemoveUserList", int32(40), userID).Return(nil).Once()

		result := DeleteUserListHandler(request(), handler)

		assert.Equal(t, okResult{nil, http.StatusNoContent}, result)

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

		result := GetUserSingleListItemHandler(request(), handler)

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

		result := GetUserSingleListItemHandler(request(), handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

		mockedListsService.AssertExpectations(t)
	})
}

func TestAddUserListItemHandler(t *testing.T) {
	request := func(useValidBody bool) *http.Request {
		var body io.Reader
		if useValidBody {
			dto := dtos.ListItemDto{
				Title: "title",
			}
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
		result := AddUserListItemHandler(request(false), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the body is valid but the insert fails", func(t *testing.T) {
		res := models.ListItem{Title: "title", ListID: int32(11)}
		mockedListsService.On("AddUserListItem", userID, &res).Return(int32(-1), &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := AddUserListItemHandler(request(true), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should add the list item if the body is valid and the insert does not fail", func(t *testing.T) {
		res := models.ListItem{Title: "title", ListID: int32(11)}
		mockedListsService.On("AddUserListItem", userID, &res).Return(int32(40), nil).Once()

		result := AddUserListItemHandler(request(true), handler)

		assert.Equal(t, okResult{int32(40), http.StatusCreated}, result)

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

		result := DeleteUserListItemHandler(request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the list item if there is no errors", func(t *testing.T) {
		mockedListsService.On("RemoveUserListItem", int32(20), int32(40), userID).Return(nil).Once()

		result := DeleteUserListItemHandler(request(), handler)

		assert.Equal(t, okResult{nil, http.StatusNoContent}, result)

		mockedListsService.AssertExpectations(t)
	})
}
