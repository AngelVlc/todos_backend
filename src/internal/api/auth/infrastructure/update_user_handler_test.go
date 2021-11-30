//+build !e2e

package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
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

	mockedRepo := authRepository.MockedAuthRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{AuthRepository: &mockedRepo, PassGen: &mockedPassGen}

	t.Run("Should return an error if the query to find the user fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update the name of the admin user", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "newAdmin"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("admin")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "It is not possible to change the admin user name")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update isAdmin of the admin user", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "admin", IsAdmin: false}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("admin")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "The admin user must be an admin")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error if the query to check if a user with the same name exists fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusR", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("ExistsUser", authDomain.UserName("wadusR")).Return(false, fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an BadRequestError if the new username already exists", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusR", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("ExistsUser", authDomain.UserName("wadusR")).Return(true, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if generate the password fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		mockedPassGen.On("GenerateFromPassword", "newPass").Return("", fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if the update fails", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusUpdated"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{Name: authDomain.UserName("wadus")}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		foundUser.Name = authDomain.UserName("wadusUpdated")
		mockedRepo.On("UpdateUser", &foundUser).Return(fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		results.CheckUnexpectedErrorResult(t, result, "Error updating the user")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the user name", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadusUpdated"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("ExistsUser", authDomain.UserName("wadusUpdated")).Return(false, nil).Once()
		foundUser2 := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		foundUser2.Name = authDomain.UserName("wadusUpdated")
		mockedRepo.On("UpdateUser", &foundUser2).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadusUpdated", userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the password", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", Password: "newPass", ConfirmPassword: "newPass"}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		mockedPassGen.On("GenerateFromPassword", "newPass").Return("hassedPass", nil).Once()
		mockedRepo.On("UpdateUser", &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the is admin value", func(t *testing.T) {
		updateReq := updateUserRequest{Name: "wadus", IsAdmin: true}
		body, _ := json.Marshal(updateReq)
		req := request(body)
		foundUser := authDomain.User{ID: int32(1), Name: authDomain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", req.Context(), int32(1)).Return(&foundUser, nil).Once()
		foundUser.IsAdmin = true
		mockedRepo.On("UpdateUser", &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), req, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.True(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})
}
