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
	"github.com/stretchr/testify/mock"
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

func TestGetUsersHandler(t *testing.T) {
	mockedUsersService := services.NewMockedUsersService()

	handler := Handler{
		usersSrv: mockedUsersService,
	}

	t.Run("Should return an errorResult if getting the users fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)

		res := []dtos.GetUsersResultDto{}
		mockedUsersService.On("GetUsers", &res).Return(&appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := GetUsersHandler(request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the users if there is no errors", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)

		res := []dtos.GetUsersResultDto{
			dtos.GetUsersResultDto{
				ID:      int32(1),
				Name:    "user1",
				IsAdmin: true,
			},
		}
		mockedUsersService.On("GetUsers", &[]dtos.GetUsersResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]dtos.GetUsersResultDto)
			*arg = res
		})

		result := GetUsersHandler(request, handler)

		assert.Equal(t, okResult{res, http.StatusOK}, result)

		mockedUsersService.AssertExpectations(t)
	})
}
