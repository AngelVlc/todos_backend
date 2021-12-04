//+build !e2e

package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	t.Run("Should return an error if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByID", request().Context(), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return the user", func(t *testing.T) {
		user := domain.User{ID: 2, Name: "user1", IsAdmin: true}
		mockedRepo.On("FindUserByID", request().Context(), int32(1)).Return(&user, nil).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(*UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(2), userRes.ID)
		assert.Equal(t, "user1", userRes.Name)
		assert.True(t, userRes.IsAdmin)
		mockedRepo.AssertExpectations(t)
	})
}
