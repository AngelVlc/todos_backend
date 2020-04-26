package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/services"
	"github.com/stretchr/testify/assert"
)

func TestAddUserHandler(t *testing.T) {
	mockedUsersService := services.NewMockedUsersService()

	handler := Handler{
		usersSrv: mockedUsersService,
	}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))

		result := AddUserHandler(request, handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if adding the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))

		mockedUsersService.On("AddUser", &dto).Return(int32(-1), &appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := AddUserHandler(request, handler)

		CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it adds the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))

		mockedUsersService.On("AddUser", &dto).Return(int32(11), nil).Once()

		result := AddUserHandler(request, handler)

		assert.Equal(t, okResult{int32(11), http.StatusCreated}, result)

		mockedUsersService.AssertExpectations(t)
	})
}
