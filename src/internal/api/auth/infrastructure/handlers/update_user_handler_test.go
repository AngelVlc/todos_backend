//go:build !e2e
// +build !e2e

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserHandlerValidations(t *testing.T) {
	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := UpdateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a create user request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

		result := UpdateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an error with a BadRequestError if the request include the passwords and they don't match", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", Password: "one", ConfirmPassword: "another"}
		body, _ := json.Marshal(updateReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := UpdateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
	})
}

func TestUpdateUserHandler(t *testing.T) {
	request := func(body []byte) *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", bytes.NewBuffer(body))
		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})
		return request
	}

	mockedUsersRepo := authRepository.MockedUsersRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{UsersRepository: &mockedUsersRepo, PassGen: &mockedPassGen}

	t.Run("Should return an error if the query to find the user fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckError(t, result, "some error")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update the name of the admin user", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "newAdmin"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("admin")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "It is not possible to change the admin user name")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update isAdmin of the admin user", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "admin", IsAdmin: false}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("admin")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "The admin user must be an admin")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("Should return an error if the query to check if a user with the same name exists fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusR", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{Name: authDomain.UserName("wadusR")}).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckError(t, result, "some error")
		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an BadRequestError if the new username already exists", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusR", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{Name: authDomain.UserName("wadusR")}).Return(&domain.User{}, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if generate the password fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedPassGen.On("GenerateFromPassword", "newPass").Return("", fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if the update fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusUpdated"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		foundUser.Name = authDomain.UserName("wadusUpdated")
		mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckUnexpectedErrorResult(t, result, "Error updating the user")
		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the user name", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusUpdated"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{Name: authDomain.UserName("wadusUpdated")}).Return(nil, nil).Once()
		foundUser2 := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		foundUser2.Name = authDomain.UserName("wadusUpdated")
		mockedUsersRepo.On("Update", req.Context(), &foundUser2).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(infrastructure.UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadusUpdated", userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the password", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		mockedPassGen.On("GenerateFromPassword", "newPass").Return("hassedPass", nil).Once()
		mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(infrastructure.UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the is admin value", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", IsAdmin: true}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedUsersRepo.On("FindUser", req.Context(), &domain.User{ID: int32(1)}).Return(&foundUser, nil).Once()
		foundUser.IsAdmin = true
		mockedUsersRepo.On("Update", req.Context(), &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(infrastructure.UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.True(t, userRes.IsAdmin)

		mockedUsersRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})
}
