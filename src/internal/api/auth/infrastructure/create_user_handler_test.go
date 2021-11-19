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
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUserHandlerValidations(t *testing.T) {
	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a create user request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request has an empty userName", func(t *testing.T) {
		createReq := createUserRequest{Name: ""}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request does not have password", func(t *testing.T) {
		createReq := createUserRequest{Name: "Wadus", Password: ""}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request passwords don't match", func(t *testing.T) {
		createReq := createUserRequest{Name: "Wadus", Password: "pass", ConfirmPassword: "othePass"}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
	})
}

func TestCreateUserHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	mockedPassGen := passgen.MockedPasswordGenerator{}
	h := handler.Handler{AuthRepository: &mockedRepo, PassGen: &mockedPassGen}

	createReq := createUserRequest{Name: "wadus", Password: "pass", ConfirmPassword: "pass", IsAdmin: true}
	body, _ := json.Marshal(createReq)

	t.Run("Should return an error if the query to check if the user exists fails", func(t *testing.T) {
		mockedRepo.On("ExistsUser", domain.UserName("wadus")).Return(false, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with a BadRequestError if a user with the same name already exist", func(t *testing.T) {
		mockedRepo.On("ExistsUser", domain.UserName("wadus")).Return(true, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if generate the password fails", func(t *testing.T) {
		mockedRepo.On("ExistsUser", domain.UserName("wadus")).Return(false, nil).Once()
		mockedPassGen.On("GenerateFromPassword", "pass").Return("", fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if create user fails", func(t *testing.T) {
		mockedRepo.On("ExistsUser", domain.UserName("wadus")).Return(false, nil).Once()
		hassedPass := "hassed"
		mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
		user := domain.User{Name: domain.UserName("wadus"), PasswordHash: hassedPass, IsAdmin: true}
		mockedRepo.On("CreateUser", &user).Return(fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating the user")
		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})

	t.Run("should create the new user", func(t *testing.T) {
		mockedRepo.On("ExistsUser", domain.UserName("wadus")).Return(false, nil).Once()
		hassedPass := "hassed"
		mockedPassGen.On("GenerateFromPassword", "pass").Return(hassedPass, nil).Once()
		user := domain.User{Name: domain.UserName("wadus"), PasswordHash: hassedPass, IsAdmin: true}
		mockedRepo.On("CreateUser", &user).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.User)
			*arg = domain.User{ID: int32(1)}
		})
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		res, isOk := okRes.Content.(UserResponse)
		require.Equal(t, true, isOk, "should be a UserResponse")
		assert.Equal(t, int32(1), res.ID)

		mockedRepo.AssertExpectations(t)
		mockedPassGen.AssertExpectations(t)
	})
}
