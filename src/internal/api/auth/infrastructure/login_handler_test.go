//+build !e2e

package infrastructure

import (
	"bytes"
	"context"
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
	"github.com/newrelic/go-agent/v3/newrelic"
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
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	loginReq := loginRequest{UserName: "wadus", Password: "pass"}
	body, _ := json.Marshal(loginReq)

	t.Run("Should return an error if the query to find the user fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(nil, fmt.Errorf("some error")).Once()

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckError(t, result, "some error")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the password does not match", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		foundUser := domain.User{PasswordHash: "hash"}
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(&foundUser, nil).Once()

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Invalid password")
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if generating the token fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.User{ID: 1, PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(&foundUser, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("", fmt.Errorf("some error")).Once()

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating jwt token")
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if generating the refresh token fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.User{ID: 1, PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(&foundUser, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("token", nil).Once()
		expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
		mockedCfgSrv.On("GetRefreshTokenExpirationDate").Return(expDate).Once()
		mockedTokenSrv.On("GenerateRefreshToken", &foundUser, expDate).Return("", fmt.Errorf("some error")).Once()

		result := LoginHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating jwt refresh token")
		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an okResult with a login response and should create the cookies although save the refresh token fails if the login is ok", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.User{ID: 1, Name: domain.UserName("user"), IsAdmin: true, PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(&foundUser, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("theToken", nil).Once()
		expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
		mockedCfgSrv.On("GetRefreshTokenExpirationDate").Return(expDate).Once()
		mockedTokenSrv.On("GenerateRefreshToken", &foundUser, expDate).Return("theRefreshToken", nil).Once()
		ctx := newrelic.NewContext(context.Background(), nil)
		mockedRepo.On("CreateRefreshTokenIfNotExist", ctx, &domain.RefreshToken{UserID: foundUser.ID, RefreshToken: "theRefreshToken", ExpirationDate: expDate}).Return(fmt.Errorf("some error")).Once()

		recorder := httptest.NewRecorder()

		mockedRepo.Wg.Add(1)
		result := LoginHandler(recorder, request, h)
		mockedRepo.Wg.Wait()

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.LoginResponse)
		require.Equal(t, true, isOk, "should be a login response")
		assert.Equal(t, "", resDto.Token)
		assert.Equal(t, "", resDto.RefreshToken)
		assert.Equal(t, int32(1), resDto.UserID)
		assert.Equal(t, "user", resDto.UserName)
		assert.True(t, resDto.IsAdmin)

		require.Equal(t, 2, len(recorder.Result().Cookies()))
		assert.Equal(t, "token", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, "theToken", recorder.Result().Cookies()[0].Value)
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[1].Name)
		assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[1].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an okResult with a login response, should create the cookies and should save the refresh token if the login is ok", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
		hashedPass := string(hashedBytes)
		foundUser := domain.User{ID: 1, Name: domain.UserName("user"), IsAdmin: true, PasswordHash: hashedPass}
		mockedRepo.On("FindUserByName", request.Context(), domain.UserName("wadus")).Return(&foundUser, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("theToken", nil).Once()
		expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
		mockedCfgSrv.On("GetRefreshTokenExpirationDate").Return(expDate).Once()
		mockedTokenSrv.On("GenerateRefreshToken", &foundUser, expDate).Return("theRefreshToken", nil).Once()
		ctx := newrelic.NewContext(context.Background(), nil)
		mockedRepo.On("CreateRefreshTokenIfNotExist", ctx, &domain.RefreshToken{UserID: foundUser.ID, RefreshToken: "theRefreshToken", ExpirationDate: expDate}).Return(nil).Once()

		recorder := httptest.NewRecorder()

		mockedRepo.Wg.Add(1)
		result := LoginHandler(recorder, request, h)
		mockedRepo.Wg.Wait()

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.LoginResponse)
		require.Equal(t, true, isOk, "should be a login response")
		assert.Equal(t, "", resDto.Token)
		assert.Equal(t, "", resDto.RefreshToken)
		assert.Equal(t, int32(1), resDto.UserID)
		assert.Equal(t, "user", resDto.UserName)
		assert.True(t, resDto.IsAdmin)

		require.Equal(t, 2, len(recorder.Result().Cookies()))
		assert.Equal(t, "token", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, "theToken", recorder.Result().Cookies()[0].Value)
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[1].Name)
		assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[1].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})
}
