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
	"github.com/stretchr/testify/assert"
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

	t.Run("Should return an errorResult with a BadRequestError if the create user request does not have userName", func(t *testing.T) {
		createReq := createUserRequest{}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request has an empty userName", func(t *testing.T) {
		s := ""
		createReq := createUserRequest{UserName: &s}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request does not have password", func(t *testing.T) {
		u := "Wadus"
		createReq := createUserRequest{UserName: &u}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request does not have password", func(t *testing.T) {
		u := "Wadus"
		p := ""
		createReq := createUserRequest{UserName: &u, Password: &p}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the create user request passwords don't match", func(t *testing.T) {
		u := "Wadus"
		p := "pass"
		confirmPass := "otherPass"
		createReq := createUserRequest{UserName: &u, Password: &p, ConfirmPassword: &confirmPass}
		body, _ := json.Marshal(createReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Passwords don't match")
	})
}

func TestCreateUserHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	mockedPassgen := authDomain.MockedPasswordGenerator{}
	h := handler.Handler{AuthRepository: &mockedRepo, PassGen: &mockedPassgen}

	u := "wadus"
	p := "pass"
	admin := true

	createReq := createUserRequest{&u, &p, &p, &admin}
	body, _ := json.Marshal(createReq)

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user name")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with a BadRequestError if a user with the same name already exist", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(&domain.AuthUser{ID: int32(1)}, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "A user with the same user name already exists")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if generate the password fails", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, nil).Once()
		pass := domain.AuthUserPassword(p)
		mockedPassgen.On("GenerateFromPassword", &pass).Return("", fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error encrypting password")
		mockedRepo.AssertExpectations(t)
		mockedPassgen.AssertExpectations(t)
	})

	t.Run("Should return an error result with an UnexpectedError if create user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, nil).Once()
		pass := domain.AuthUserPassword(p)
		hassedPass := "hassed"
		mockedPassgen.On("GenerateFromPassword", &pass).Return(hassedPass, nil).Once()
		user := domain.AuthUser{Name: domain.AuthUserName(u), PasswordHash: hassedPass, IsAdmin: true}
		mockedRepo.On("CreateUser", &user).Return(nil, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating the user")
		mockedRepo.AssertExpectations(t)
		mockedPassgen.AssertExpectations(t)
	})

	t.Run("should create the new user", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, nil).Once()
		pass := domain.AuthUserPassword(p)
		hassedPass := "hassed"
		mockedPassgen.On("GenerateFromPassword", &pass).Return(hassedPass, nil).Once()
		user := domain.AuthUser{Name: domain.AuthUserName(u), PasswordHash: hassedPass, IsAdmin: true}
		mockedRepo.On("CreateUser", &user).Return(int32(1), nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := CreateUserHandler(httptest.NewRecorder(), request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		res, isOk := okRes.Content.(int32)
		require.Equal(t, true, isOk, "should be a token response")
		assert.Equal(t, int32(1), res)

		mockedRepo.AssertExpectations(t)
		mockedPassgen.AssertExpectations(t)
	})
}
