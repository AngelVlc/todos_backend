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
	"github.com/stretchr/testify/mock"
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

	matchFn := func(i *int32) bool {
		return *i == 1
	}

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(nil, fmt.Errorf("some error")).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by id")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError if the user does not exits", func(t *testing.T) {
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(nil, nil).Once()

		result := GetUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "The user does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return the user", func(t *testing.T) {
		user := domain.User{ID: 2, Name: "user1", IsAdmin: true}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&user, nil).Once()

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
