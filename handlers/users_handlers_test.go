package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddUserHandler(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		return request
	}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		result := AddUserHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if adding the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("AddUser", &dto).Return(int32(-1), &appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := AddUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), handler)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it adds the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("AddUser", &dto).Return(int32(11), nil).Once()

		result := AddUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), handler)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(11), okRes.content)

		mockedUsersService.AssertExpectations(t)
	})
}

func TestGetUsersHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		return request
	}

	t.Run("Should return an errorResult if getting the users fails", func(t *testing.T) {
		mockedUsersService.On("GetUsers").Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUsersHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the users if there is no errors", func(t *testing.T) {
		res := []*dtos.UserResponseDto{
			{
				ID:      int32(1),
				Name:    "user1",
				IsAdmin: true,
			},
		}
		mockedUsersService.On("GetUsers").Return(res, nil).Once()

		result := GetUsersHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.content.([]*dtos.UserResponseDto)
		require.Equal(t, true, isOk, "should be a user result dto")
		require.Equal(t, len(res), len(resDto))
		assert.Equal(t, res[0].ID, resDto[0].ID)
		assert.Equal(t, res[0].Name, resDto[0].Name)
		assert.Equal(t, res[0].IsAdmin, resDto[0].IsAdmin)

		mockedUsersService.AssertExpectations(t)
	})
}

func TestDeleteUserHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		return request
	}

	t.Run("Should return an errorResult if getting the user lists fails", func(t *testing.T) {
		mockedListsService.On("GetUserLists", int32(40), &[]dtos.GetListsResultDto{}).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the user has some list", func(t *testing.T) {
		res := []dtos.GetListsResultDto{
			dtos.GetListsResultDto{
				ID:   int32(1),
				Name: "list1",
			},
		}
		mockedListsService.On("GetUserLists", int32(40), &[]dtos.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]dtos.GetListsResultDto)
			*arg = res
		})

		result := DeleteUserHandler(httptest.NewRecorder(), request(), handler)

		CheckBadRequestErrorResult(t, result, "The user has lists")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		res := []dtos.GetListsResultDto{}
		mockedListsService.On("GetUserLists", int32(40), &[]dtos.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]dtos.GetListsResultDto)
			*arg = res
		})
		mockedUsersService.On("RemoveUser", int32(40)).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the user if there is no errors", func(t *testing.T) {
		res := []dtos.GetListsResultDto{}
		mockedListsService.On("GetUserLists", int32(40), &[]dtos.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]dtos.GetListsResultDto)
			*arg = res
		})
		mockedUsersService.On("RemoveUser", int32(40)).Return(nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, okRes.content)

		mockedUsersService.AssertExpectations(t)
		mockedListsService.AssertExpectations(t)
	})
}

func TestUpdateUserHandler(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		return request
	}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		result := UpdateUserHandler(httptest.NewRecorder(), request(nil), handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if updating the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(nil, &appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), handler)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it updates the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		user := models.User{Name: "updated"}

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(&user, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), handler)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		resDto, isOk := okRes.content.(*models.User)
		require.Equal(t, true, isOk, "should be a list result dto")
		assert.Equal(t, user.Name, resDto.Name)

		mockedUsersService.AssertExpectations(t)
	})
}

func TestGetUserHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})

		return request
	}

	t.Run("Should return an errorResult if updating the user fails", func(t *testing.T) {
		mockedUsersService.On("FindUserByID", int32(40)).Return(nil, &appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), handler)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the user info", func(t *testing.T) {
		user := models.User{
			Name:    "user",
			IsAdmin: true,
			ID:      40,
		}

		mockedUsersService.On("FindUserByID", int32(40)).Return(&user, nil).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), handler)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		resDto, isOk := okRes.content.(dtos.UserResponseDto)
		require.Equal(t, true, isOk, "should be a user result dto")
		assert.Equal(t, user.ID, resDto.ID)
		assert.Equal(t, user.Name, resDto.Name)
		assert.Equal(t, user.IsAdmin, resDto.IsAdmin)

		mockedUsersService.AssertExpectations(t)
	})
}
