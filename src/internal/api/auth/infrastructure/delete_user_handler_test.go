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
)

func TestDeleteUserHandler(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedRepo := authRepository.MockedAuthRepository{}
	h := handler.Handler{AuthRepository: &mockedRepo}

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByID", int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by id")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError if the user does not exits", func(t *testing.T) {
		mockedRepo.On("FindUserByID", int32(1)).Return(nil, nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "The user does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError when deleting the admin user", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("admin")}
		mockedRepo.On("FindUserByID", int32(1)).Return(&foundUser, nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "It is not possible to delete the admin user")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the delete fails", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedRepo.On("FindUserByID", int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("DeleteUser", int32(1)).Return(fmt.Errorf("some error")).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error deleting the user")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should delete the user", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedRepo.On("FindUserByID", int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("DeleteUser", int32(1)).Return(nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckOkResult(t, result, http.StatusNoContent)
		mockedRepo.AssertExpectations(t)
	})
}
