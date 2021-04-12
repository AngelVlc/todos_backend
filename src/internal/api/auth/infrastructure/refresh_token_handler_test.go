//+build !e2e

package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenHandlerValidations(t *testing.T) {
	h := handler.Handler{}

	t.Run("Should return an bad request error if the request does not come with the refresh token cookie", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckBadRequestErrorResult(t, result, "Missing refresh token cookie")
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	mockedRepo := authRepository.MockedAuthRepository{}
	mockedCfgSrv := sharedApp.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token is not valid", func(t *testing.T) {
		mockedTokenSrv.On("ParseToken", "badToken").Return(nil, fmt.Errorf("some error")).Once()

		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("badToken"))

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "Invalid refresh token")
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if getting the user by id fails", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		mockedRepo.On("FindUserByID", int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user id")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the user no longer exists", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		mockedRepo.On("FindUserByID", int32(1)).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "The user no longer exists")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if getting the refresh token fails", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		mockedRepo.On("FindUserByID", int32(1)).Return(&domain.User{}, nil).Once()
		mockedRepo.On("FindRefreshTokenForUser", "token", int32(1)).Return(nil, fmt.Errorf("some error")).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting the refresh token")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token does not exist", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		mockedRepo.On("FindUserByID", int32(1)).Return(&domain.User{}, nil).Once()
		mockedRepo.On("FindRefreshTokenForUser", "token", int32(1)).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "The refresh token is not valid")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnexpectedError if generate the new token fails", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		foundUser := domain.User{}
		mockedRepo.On("FindUserByID", int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("FindRefreshTokenForUser", "token", int32(1)).Return(&domain.RefreshToken{}, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("", fmt.Errorf("some error")).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error creating jwt token")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the token and should create the cookie if the refresh token is valid", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
		rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
		mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("token"))
		foundUser := domain.User{}
		mockedRepo.On("FindUserByID", int32(1)).Return(&foundUser, nil).Once()
		mockedRepo.On("FindRefreshTokenForUser", "token", int32(1)).Return(&domain.RefreshToken{}, nil).Once()
		mockedTokenSrv.On("GenerateToken", &foundUser).Return("token", nil).Once()

		recorder := httptest.NewRecorder()
		result := RefreshTokenHandler(recorder, request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.TokenResponse)
		require.Equal(t, true, isOk, "should be a token response")
		assert.Equal(t, "token", resDto.Token)
		assert.Equal(t, "", resDto.RefreshToken)

		require.Equal(t, 1, len(recorder.Result().Cookies()))
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, resDto.RefreshToken, recorder.Result().Cookies()[0].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
		mockedTokenSrv.AssertExpectations(t)
	})
}
