package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/services"
	"github.com/stretchr/testify/assert"
)

var (
	mockedAuthService  = services.NewMockedAuthService()
	mockedUsersService = services.NewMockedUsersService()
	mockedListsService = services.NewMockedListsService()

	handler = Handler{
		usersSrv: mockedUsersService,
		authSrv:  mockedAuthService,
		listsSrv: mockedListsService,
	}
)

func TestTokenHandler(t *testing.T) {
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
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		mockedUsersService.On("FindUserByName", login.UserName).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with a bad request error if the user does not exist", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		mockedUsersService.On("FindUserByName", login.UserName).Return(nil, nil).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "The user does not exist")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with a bad request error if the user exists and the password is not valid", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		user := models.User{}

		mockedUsersService.On("FindUserByName", login.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, login.Password).Return(fmt.Errorf("some error")).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Invalid password")

		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an error result with an unexpected error if the user and the password are correct but the tokens generation fails", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		user := models.User{}

		mockedUsersService.On("FindUserByName", login.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, login.Password).Return(nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		result := TokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedUsersService.AssertExpectations(t)
		mockedAuthService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the tokens if the user and the password are correct", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		user := models.User{}

		tokens := map[string]string{
			"token":        "theToken",
			"refreshToken": "theRefreshToken",
		}

		mockedUsersService.On("FindUserByName", login.UserName).Return(&user, nil).Once()
		mockedUsersService.On("CheckIfUserPasswordIsOk", &user, login.Password).Return(nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(tokens, nil).Once()

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		result := TokenHandler(recorder, request, handler)

		assert.Equal(t, okResult{tokens, http.StatusOK}, result)
		assertSuccessResponse(t, recorder)

		mockedUsersService.AssertExpectations(t)
		mockedAuthService.AssertExpectations(t)
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	t.Run("Should return an errorResult with a BadRequestError if the body is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", strings.NewReader("wadus"))

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "Invalid body")
	})

	t.Run("Should return an errorResult with a BadRequestError if the body does not have refresh token", func(t *testing.T) {
		login := struct{}{}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "RefreshToken is mandatory")
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token is not valid", func(t *testing.T) {
		refreshToken := models.RefreshToken{
			RefreshToken: "theRefreshToken",
		}
		body, _ := json.Marshal(refreshToken)

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))

		mockedAuthService.On("ParseRefreshToken", refreshToken.RefreshToken).Return(nil, &appErrors.UnauthorizedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnauthorizedErrorErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the refresh token is valid but getting the user by id fails", func(t *testing.T) {
		refreshToken := models.RefreshToken{
			RefreshToken: "theRefreshToken",
		}
		body, _ := json.Marshal(refreshToken)

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 1,
		}

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))

		mockedAuthService.On("ParseRefreshToken", refreshToken.RefreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with a BadRequestError if the refresh token is valid but the user does not exist", func(t *testing.T) {
		refreshToken := models.RefreshToken{
			RefreshToken: "theRefreshToken",
		}
		body, _ := json.Marshal(refreshToken)

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 1,
		}

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))

		mockedAuthService.On("ParseRefreshToken", refreshToken.RefreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckBadRequestErrorResult(t, result, "The user is no longer valid")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if the refresh token is valid, the user exists but getting the tokens fails", func(t *testing.T) {
		refreshToken := models.RefreshToken{
			RefreshToken: "theRefreshToken",
		}
		body, _ := json.Marshal(refreshToken)

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 1,
		}

		user := models.User{}

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))

		mockedAuthService.On("ParseRefreshToken", refreshToken.RefreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(&user, nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(nil, &appErrors.UnexpectedError{Msg: "Some error"}).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, handler)

		CheckUnexpectedErrorResult(t, result, "Some error")

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})

	t.Run("Should return an ok result with the tokens if refresh token is valid and the users exists", func(t *testing.T) {
		refreshToken := models.RefreshToken{
			RefreshToken: "theRefreshToken",
		}
		body, _ := json.Marshal(refreshToken)

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 1,
		}

		user := models.User{}

		tokens := map[string]string{
			"token":        "theToken",
			"refreshToken": "theRefreshToken",
		}

		request, _ := http.NewRequest(http.MethodPost, "/wadus", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		mockedAuthService.On("ParseRefreshToken", refreshToken.RefreshToken).Return(&rtInfo, nil)
		mockedUsersService.On("FindUserByID", rtInfo.UserID).Return(&user, nil).Once()
		mockedAuthService.On("GetTokens", &user).Return(tokens, nil).Once()

		result := RefreshTokenHandler(recorder, request, handler)

		assert.Equal(t, okResult{tokens, http.StatusOK}, result)
		assertSuccessResponse(t, recorder)

		mockedAuthService.AssertExpectations(t)
		mockedUsersService.AssertExpectations(t)
	})
}

func assertSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	assert.Equal(t, 1, len(recorder.Result().Cookies()))
	assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
	assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[0].Value)
	assert.True(t, recorder.Result().Cookies()[0].HttpOnly)
}
