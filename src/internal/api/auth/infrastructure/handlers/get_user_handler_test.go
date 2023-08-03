//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserHandler_Returns_An_Error_If_The_Query_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedRepo}

	mockedRepo.On("FindUser", request().Context(), domain.UserEntity{ID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

	result := GetUserHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedRepo.AssertExpectations(t)
}

func TestGetUserHandler_Returns_The_User(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedRepo}

	nvo, _ := domain.NewUserNameValueObject("user1")
	user := domain.UserEntity{ID: 2, Name: nvo, IsAdmin: true}
	mockedRepo.On("FindUser", request().Context(), domain.UserEntity{ID: int32(1)}).Return(&user, nil).Once()

	result := GetUserHandler(httptest.NewRecorder(), request(), h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	userRes, isOk := okRes.Content.(*infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a user response")

	assert.Equal(t, int32(2), userRes.ID)
	assert.Equal(t, "user1", userRes.Name)
	assert.True(t, userRes.IsAdmin)
	mockedRepo.AssertExpectations(t)
}
