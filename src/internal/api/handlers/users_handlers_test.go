//+build !e2e

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/models"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddUserHandler(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		return request
	}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		result := AddUserHandler(httptest.NewRecorder(), request(strings.NewReader("wadus")), h)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if adding the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("AddUser", &dto).Return(int32(-1), &appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := AddUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it adds the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("AddUser", &dto).Return(int32(11), nil).Once()

		result := AddUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		okRes := CheckOkResult(t, result, http.StatusCreated)
		assert.Equal(t, int32(11), okRes.Content)

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

		result := GetUsersHandler(httptest.NewRecorder(), request(), h)

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

		result := GetUsersHandler(httptest.NewRecorder(), request(), h)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.([]*dtos.UserResponseDto)
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
		mockedListsService.On("GetUserLists", int32(40)).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the user has some list", func(t *testing.T) {
		res := []*dtos.ListResponseDto{
			{
				ID:   int32(1),
				Name: "list1",
			},
		}
		mockedListsService.On("GetUserLists", int32(40)).Return(res, nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		CheckBadRequestErrorResult(t, result, "The user has lists")

		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult if the delete fails", func(t *testing.T) {
		mockedListsService.On("GetUserLists", int32(40)).Return([]*dtos.ListResponseDto{}, nil).Once()
		mockedUsersService.On("RemoveUser", int32(40)).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
		mockedListsService.AssertExpectations(t)
	})

	t.Run("Should delete the user if there is no errors", func(t *testing.T) {
		mockedListsService.On("GetUserLists", int32(40)).Return([]*dtos.ListResponseDto{}, nil).Once()
		mockedUsersService.On("RemoveUser", int32(40)).Return(nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		noContentRes := CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, noContentRes.Content)

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
		result := UpdateUserHandler(httptest.NewRecorder(), request(nil), h)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if updating the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(&appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it updates the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		noContentRes := CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, noContentRes.Content)

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

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

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

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

		okRes := CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(dtos.UserResponseDto)
		require.Equal(t, true, isOk, "should be a user result dto")
		assert.Equal(t, user.ID, resDto.ID)
		assert.Equal(t, user.Name, resDto.Name)
		assert.Equal(t, user.IsAdmin, resDto.IsAdmin)

		mockedUsersService.AssertExpectations(t)
	})
}
