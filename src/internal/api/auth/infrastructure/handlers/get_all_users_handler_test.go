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
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllUsersHandler_Returns_An_UnexpectedError_If_The_Query_Fails(t *testing.T) {
	mockedRepo := authRepository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedRepo}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mockedRepo.On("GetAll", request.Context()).Return(nil, fmt.Errorf("some error")).Once()

	result := GetAllUsersHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting users")
	mockedRepo.AssertExpectations(t)
}

func TestGetAllUsersHandler_Returns_The_Users(t *testing.T) {
	mockedRepo := authRepository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedRepo}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	found := []domain.User{
		{ID: 2, Name: "user1", IsAdmin: true},
		{ID: 5, Name: "user2", IsAdmin: false},
	}
	mockedRepo.On("GetAll", request.Context()).Return(found, nil)
	result := GetAllUsersHandler(httptest.NewRecorder(), request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	userRes, isOk := okRes.Content.([]*infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be an array of user response")

	require.Equal(t, len(userRes), 2)
	assert.Equal(t, int32(2), userRes[0].ID)
	assert.Equal(t, "user1", userRes[0].Name)
	assert.True(t, userRes[0].IsAdmin)
	assert.Equal(t, int32(5), userRes[1].ID)
	assert.Equal(t, "user2", userRes[1].Name)
	assert.False(t, userRes[1].IsAdmin)

	mockedRepo.AssertExpectations(t)
}
