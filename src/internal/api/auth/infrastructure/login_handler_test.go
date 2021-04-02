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
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandlerValidations(t *testing.T) {
	h := handler.Handler{}

	t.Run("Should return an errorResult with a BadRequestError if the request does not have body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body is not a login request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("wadus"))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the login request has an empty userName", func(t *testing.T) {
		loginReq := loginRequest{UserName: ""}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the login request does not have password", func(t *testing.T) {
		loginReq := loginRequest{UserName: "wadus", Password: ""}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password can not be empty")
	})
}

func TestLoginHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	mockedCfgSrv := sharedApp.MockedConfigurationService{}
	h := handler.Handler{AuthRepository: &mockedRepo, CfgSrv: &mockedCfgSrv}

	loginReq := loginRequest{UserName: "wadus", Password: "pass"}
	body, _ := json.Marshal(loginReq)

	t.Run("Should return an errorResult with an UnexpectedError if the query to find the user fails", func(t *testing.T) {
		mockedRepo.On("FindUserByName", domain.UserName("wadus")).Return(nil, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user name")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the user does not exist", func(t *testing.T) {
		mockedRepo.On("FindUserByName", domain.UserName("wadus")).Return(nil, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "The user does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the password does not match", func(t *testing.T) {
		foundUser := domain.User{PasswordHash: "hash"}
		mockedRepo.On("FindUserByName", domain.UserName("wadus")).Return(&foundUser, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid password")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the tokens and should create the cookie if the login is correct", func(t *testing.T) {
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.User{PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", domain.UserName("wadus")).Return(&foundUser, nil).Once()
		expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
		mockedCfgSrv.On("GetTokenExpirationDate").Return(expDate).Once()
		mockedCfgSrv.On("GetRefreshTokenExpirationDate").Return(expDate).Once()
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(2)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		recorder := httptest.NewRecorder()
		result := LoginHandler(recorder, request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.TokenResponse)
		require.Equal(t, true, isOk, "should be a token response")
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTc0NzY0MDAsImlzQWRtaW4iOmZhbHNlLCJ1c2VySWQiOjAsInVzZXJOYW1lIjoiIn0.vZjb1EWpNfdjeR1roJHhRnFKsPXIKMPZlgxigdupIHo", resDto.Token)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTc0NzY0MDAsInVzZXJJZCI6MH0.9E8npy60pAIzzvv7V0The5457bVcrMAxzbdYPo63kMo", resDto.RefreshToken)

		require.Equal(t, 1, len(recorder.Result().Cookies()))
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, resDto.RefreshToken, recorder.Result().Cookies()[0].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
	})
}
