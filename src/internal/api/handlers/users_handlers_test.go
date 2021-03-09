//+build !e2e

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/services"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserHandler(t *testing.T) {
	request := func(body io.Reader) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", body)
		request = mux.SetURLVars(request, map[string]string{
			"id": "40",
		})
		return request
	}

	mockedUsersService := services.NewMockedUsersService()
	h := handler.Handler{UsersSrv: mockedUsersService}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		result := UpdateUserHandler(httptest.NewRecorder(), request(nil), h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult if updating the user fails", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(&appErrors.BadRequestError{Msg: "Some error"}).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		results.CheckBadRequestErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an okResult when it updates the user", func(t *testing.T) {
		dto := dtos.UserDto{}
		body, _ := json.Marshal(dto)

		mockedUsersService.On("UpdateUser", int32(40), &dto).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(bytes.NewBuffer(body)), h)

		noContentRes := results.CheckOkResult(t, result, http.StatusNoContent)
		assert.Nil(t, noContentRes.Content)

		mockedUsersService.AssertExpectations(t)
	})
}
