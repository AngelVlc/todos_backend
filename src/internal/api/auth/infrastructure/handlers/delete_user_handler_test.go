//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
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

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo}

	t.Run("Should return an error if the query to find the user fails", func(t *testing.T) {
		mockedUsersRepo.On("FindUser", request().Context(), &domain.User{ID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckError(t, result, "some error")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an errorResult with a BadRequestError when deleting the admin user", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("admin")}
		mockedUsersRepo.On("FindUser", request().Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckBadRequestErrorResult(t, result, "It is not possible to delete the admin user")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the delete fails", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", request().Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedUsersRepo.On("Delete", request().Context(), &domain.User{ID: int32(1)}).Return(fmt.Errorf("some error")).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckUnexpectedErrorResult(t, result, "Error deleting the user")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("Should delete the user", func(t *testing.T) {
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", request().Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedUsersRepo.On("Delete", request().Context(), &domain.User{ID: int32(1)}).Return(nil).Once()

		result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

		results.CheckOkResult(t, result, http.StatusNoContent)
		mockedUsersRepo.AssertExpectations(t)
	})
}
