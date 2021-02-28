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

	t.Run("Should return an errorResult with a BadRequestError if the login request does not have userName", func(t *testing.T) {
		loginReq := loginRequest{}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the login request has an empty userName", func(t *testing.T) {
		s := ""
		loginReq := loginRequest{UserName: &s}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "UserName can not be empty")
	})

	t.Run("Should return an errorResult with a BadRequestError if the login request does not have password", func(t *testing.T) {
		u := "Wadus"
		loginReq := loginRequest{UserName: &u}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the login request does not have password", func(t *testing.T) {
		u := "Wadus"
		p := ""
		loginReq := loginRequest{UserName: &u, Password: &p}
		body, _ := json.Marshal(loginReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Password can not be empty")
	})
}

func TestLoginHandler(t *testing.T) {
	mockedRepo := MockedAuthRepository{}
	mockedCfgSrv := sharedApp.MockedConfigurationService{}
	h := handler.Handler{AuthRepository: &mockedRepo, CfgSrv: &mockedCfgSrv}

	u := "wadus"
	p := "pass"
	loginReq := loginRequest{&u, &p}
	body, _ := json.Marshal(loginReq)

	t.Run("Should return an errorResult with an UnexpectedError if the user does not exist", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, fmt.Errorf("some error")).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user name")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the user does not exist", func(t *testing.T) {
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(nil, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "The user does not exist")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the password does not match", func(t *testing.T) {
		foundUser := domain.AuthUser{PasswordHash: "hash"}
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(&foundUser, nil).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid password")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the tokens and should create the cookie if the login is correct", func(t *testing.T) {
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.AuthUser{PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", (*domain.AuthUserName)(&u)).Return(&foundUser, nil).Once()
		mockedCfgSrv.On("TokenExpirationInSeconds").Return(time.Second * 30).Once()
		mockedCfgSrv.On("RefreshTokenExpirationInSeconds").Return(time.Hour).Once()
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(2)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		recorder := httptest.NewRecorder()
		result := LoginHandler(recorder, request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.TokenResponse)
		require.Equal(t, true, isOk, "should be a token response")
		assert.Equal(t, 160, len(resDto.Token))
		assert.Equal(t, 120, len(resDto.RefreshToken))

		require.Equal(t, 1, len(recorder.Result().Cookies()))
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, resDto.RefreshToken, recorder.Result().Cookies()[0].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
	})
}
