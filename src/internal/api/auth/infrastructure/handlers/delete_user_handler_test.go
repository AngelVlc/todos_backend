//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
)

func TestDeleteUserHandler_Returns_An_Error_If_The_Query_To_Find_The_User_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo}

	mockedUsersRepo.On("FindUser", request().Context(), domain.UserRecord{ID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

	results.CheckError(t, result, "some error")
	mockedUsersRepo.AssertExpectations(t)
}

func TestDeleteUserHandler_Returns_An_ErrorResult_With_A_BadRequestError_When_Deleting_The_Admin_User(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo}

	foundUser := domain.UserRecord{Name: "admin"}
	mockedUsersRepo.On("FindUser", request().Context(), domain.UserRecord{ID: 1}).Return(&foundUser, nil).Once()

	result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

	results.CheckBadRequestErrorResult(t, result, "It is not possible to delete the admin user")
	mockedUsersRepo.AssertExpectations(t)
}

func TestDeleteUserHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_The_Delete_Fails(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo}

	foundUser := domain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", request().Context(), domain.UserRecord{ID: 1}).Return(&foundUser, nil).Once()
	mockedUsersRepo.On("Delete", request().Context(), domain.UserRecord{ID: 1}).Return(fmt.Errorf("some error")).Once()

	result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

	results.CheckUnexpectedErrorResult(t, result, "Error deleting the user")
	mockedUsersRepo.AssertExpectations(t)
}

func TestDeleteUserHandler_Deletes_The_User(t *testing.T) {
	request := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := repository.MockedUsersRepository{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo}

	foundUser := domain.UserRecord{Name: "wadus"}
	mockedUsersRepo.On("FindUser", request().Context(), domain.UserRecord{ID: 1}).Return(&foundUser, nil).Once()
	mockedUsersRepo.On("Delete", request().Context(), domain.UserRecord{ID: 1}).Return(nil).Once()

	result := DeleteUserHandler(httptest.NewRecorder(), request(), h)

	results.CheckOkResult(t, result, http.StatusNoContent)
	mockedUsersRepo.AssertExpectations(t)
}
