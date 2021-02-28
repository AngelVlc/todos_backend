//+build !e2e

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AngelVlc/todos/internal/api/dtos"
	appErrors "github.com/AngelVlc/todos/internal/api/errors"
	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/AngelVlc/todos/internal/api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mockedAuthService      = services.NewMockedAuthService()
	mockedUsersService     = services.NewMockedUsersService()
	mockedListsService     = services.NewMockedListsService()
	mockedListItemsService = services.NewMockedListItemsService()

	handler = Handler{
		usersSrv:     mockedUsersService,
		authSrv:      mockedAuthService,
		listsSrv:     mockedListsService,
		listItemsSrv: mockedListItemsService,
	}
)

func TestTokenHandler(t *testing.T) {
	validLogin := dtos.TokenDto{UserName: "wadus", Password: "pass"}

	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body does not have user name", func(t *testing.T) {
		login := struct {
			Password string
		}{
			"pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "UserName is mandatory")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body does not have password", func(t *testing.T) {
		login := struct {
			UserName string
		}{
			"wadus",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Password is mandatory")
	})

	t.Run("Should return an error result with an unexpexted error if getting the user fails", func(t *testing.T) {
		body, _ := json.Marshal(validLogin)

		mockedUsersService.On("FindUserByName", validLogin.UserName).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with a bad request error if the user does not exist", func(t *testing.T) {
		body, _ := json.Marshal(validLogin)

		mockedUsersService.On("FindUserByName", validLogin.UserName).Return(nil, nil).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "The user does not exist")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with a bad request error if the user exists and the password is not valid", func(t *testing.T) {
		body, _ := json.Marshal(validLogin)

		user := models.User{}

		mockedUsersService.On("FindUserByName", validLogin.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, validLogin.Password).Return(fmt.Errorf("some error")).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Invalid password")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with an unexpected error if the user and the password are correct but the tokens generation fails", func(t *testing.T) {
		body, _ := json.Marshal(validLogin)

		user := models.User{}

		mockedUsersService.On("FindUserByName", validLogin.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, validLogin.Password).Return(nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
		mockedAuthService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the tokens if the user and the password are correct", func(t *testing.T) {
		body, _ := json.Marshal(validLogin)

		user := models.User{}

		tokens := newTokenResultDto()

		mockedUsersService.On("FindUserByName", validLogin.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, validLogin.Password).Return(nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(tokens, nil).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		result := TokenHandler(recorder, request, handler)

		checkTokensResponse(t, result, tokens)
		checkResponseCookie(t, recorder)

		mockedUsersService.AssertExpectations(t)
		mockedAuthService.AssertExpectations(t)
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	refreshToken := "theRefreshToken"

	rtInfo := models.RefreshTokenClaimsInfo{UserID: 1}

	getRefreshTokenCookie := func() *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: refreshToken}
	}

	t.Run("Should return an errorResult with a BadRequestError if there isn't refresh token cookie", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Missing refresh token cookie")
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/wadus", nil)
		request.AddCookie(getRefreshTokenCookie())

		mockedAuthService.On("ParseRefreshToken", refreshToken).Return(nil, &appErrors.UnauthorizedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnauthorizedErrorErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the refresh token is valid but getting the user by id fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/wadus", nil)
		request.AddCookie(getRefreshTokenCookie())

		mockedAuthService.On("ParseRefreshToken", refreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the refresh token is valid but the user does not exist", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/wadus", nil)
		request.AddCookie(getRefreshTokenCookie())

		mockedAuthService.On("ParseRefreshToken", refreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "The user is no longer valid")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the refresh token is valid, the user exists but getting the tokens fails", func(t *testing.T) {
		user := models.User{}

		request, _ := http.NewRequest(http.MethodPost, "/wadus", nil)
		request.AddCookie(getRefreshTokenCookie())

		mockedAuthService.On("ParseRefreshToken", refreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(&user, nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the tokens if refresh token is valid and the users exists", func(t *testing.T) {
		user := models.User{}

		tokens := newTokenResultDto()

		request, _ := http.NewRequest(http.MethodPost, "/wadus", nil)
		request.AddCookie(getRefreshTokenCookie())
		recorder := httptest.NewRecorder()

		mockedAuthService.On("ParseRefreshToken", refreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(&user, nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(tokens, nil).Once()

		result := RefreshTokenHandler(recorder, request, handler)

		checkTokensResponse(t, result, tokens)
		checkResponseCookie(t, recorder)

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})
}

func newTokenResultDto() *dtos.TokenResponseDto {
	return &dtos.TokenResponseDto{Token: "theToken", RefreshToken: "theRefreshToken"}
}

func checkTokensResponse(t *testing.T, result HandlerResult, expectedTokens *dtos.TokenResponseDto) {
	okRes := CheckOkResult(t, result, http.StatusOK)
	tokenDto, isTokenResultDto := okRes.content.(*dtos.TokenResponseDto)
	require.Equal(t, true, isTokenResultDto, "should be a token result dto")
	assert.Equal(t, expectedTokens.Token, tokenDto.Token)
	assert.Equal(t, expectedTokens.RefreshToken, tokenDto.RefreshToken)
}

func checkResponseCookie(t *testing.T, recorder *httptest.ResponseRecorder) {
	assert.Equal(t, 1, len(recorder.Result().Cookies()))
	assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
	assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[0].Value)
	assert.True(t, recorder.Result().Cookies()[0].HttpOnly)
}
