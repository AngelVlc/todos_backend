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

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		pass := "one"
		confirmPass := "another"
		updateReq := updateUserRequest{Password: &pass, ConfirmPassword: &confirmPass}
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
	mockedPassGen := authDomain.MockedPasswordGenerator{}
	h := handler.Handler{AuthRepository: &mockedRepo, PassGen: &mockedPassGen}

	matchFn := func(i *int32) bool {
		return *i == 1
	}

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		updateReq := updateUserRequest{}
		body, _ := json.Marshal(updateReq)
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(nil, fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by id")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the user does not exist", func(t *testing.T) {
		updateReq := updateUserRequest{}
		body, _ := json.Marshal(updateReq)
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(nil, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The user does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update the name of the admin user", func(t *testing.T) {
		name := "newAdmin"
		updateReq := updateUserRequest{UserName: &name}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{Name: domain.UserName("admin")}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "It is not possible to change the admin user name")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if tries to update the name of the admin user", func(t *testing.T) {
		name := "admin"
		isAdmin := false
		updateReq := updateUserRequest{UserName: &name, IsAdmin: &isAdmin}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{Name: domain.UserName("admin")}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckBadRequestErrorResult(t, result, "The admin user must be an admin")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if generate the password fails", func(t *testing.T) {
		name := "wadus"
		pass := "newPass"
		updateReq := updateUserRequest{UserName: &name, Password: &pass, ConfirmPassword: &pass}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()
		authPass := domain.UserPassword(pass)
		mockedPassGen.On("GenerateFromPassword", &authPass).Return("", fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if the update fails", func(t *testing.T) {
		name := "wadusUpdated"
		updateReq := updateUserRequest{UserName: &name}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{Name: domain.UserName("wadus")}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()
		foundUser.Name = authDomain.UserName(name)
		mockedRepo.On("UpdateUser", &foundUser).Return(fmt.Errorf("some error")).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		results.CheckUnexpectedErrorResult(t, result, "Error updating the user")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the user name", func(t *testing.T) {
		name := "wadusUpdated"
		updateReq := updateUserRequest{UserName: &name}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{ID: int32(1), Name: domain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()
		foundUser.Name = authDomain.UserName(name)
		mockedRepo.On("UpdateUser", &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(*UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, name, userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the password", func(t *testing.T) {
		newPass := "newPass"
		updateReq := updateUserRequest{Password: &newPass, ConfirmPassword: &newPass}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{ID: int32(1), Name: domain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()
		authPass := domain.UserPassword(newPass)
		mockedPassGen.On("GenerateFromPassword", &authPass).Return("hassedPass", nil).Once()
		mockedRepo.On("UpdateUser", &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(*UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.False(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should update the is admin value", func(t *testing.T) {
		isAdmin := true
		updateReq := updateUserRequest{IsAdmin: &isAdmin}
		body, _ := json.Marshal(updateReq)
		foundUser := domain.User{ID: int32(1), Name: domain.UserName("wadus"), IsAdmin: false}
		mockedRepo.On("FindUserByID", mock.MatchedBy(matchFn)).Return(&foundUser, nil).Once()
		foundUser.IsAdmin = true
		mockedRepo.On("UpdateUser", &foundUser).Return(nil).Once()

		result := UpdateUserHandler(httptest.NewRecorder(), request(body), h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		userRes, isOk := okRes.Content.(*UserResponse)
		require.Equal(t, true, isOk, "should be a user response")

		assert.Equal(t, int32(1), userRes.ID)
		assert.Equal(t, "wadus", userRes.Name)
		assert.True(t, userRes.IsAdmin)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})
}
