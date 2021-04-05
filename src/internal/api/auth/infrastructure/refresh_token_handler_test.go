//+build !e2e

package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
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
	h := handler.Handler{AuthRepository: &mockedRepo, CfgSrv: &mockedCfgSrv}

	cfgSvc := sharedApp.NewRealConfigurationService(sharedApp.NewOsEnvGetter())
	tokenSvc := domain.NewTokenService(cfgSvc)
	refreshTokenExpDate := cfgSvc.GetRefreshTokenExpirationDate()
	refreshToken, _ := tokenSvc.GenerateRefreshToken(&domain.User{ID: 1}, refreshTokenExpDate)

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("badToken"))

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "Error parsing the refresh token")
	})

	t.Run("Should return an errorResult with an UnexpectedError if getting the user by id fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return(cfgSvc.GetJwtSecret()).Once()
		authUser := domain.User{ID: 1}
		request.AddCookie(getRefreshTokenCookie(refreshToken))
		mockedRepo.On("FindUserByID", authUser.ID).Return(nil, fmt.Errorf("some error")).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user id")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the user no longer exists", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return(cfgSvc.GetJwtSecret()).Once()
		authUser := domain.User{ID: 1}
		request.AddCookie(getRefreshTokenCookie(refreshToken))
		mockedRepo.On("FindUserByID", authUser.ID).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "The user no longer exists")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the tokens and should create the cookie if the refresh token is valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return(os.Getenv("JWT_SECRET")).Times(3)
		expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
		mockedCfgSrv.On("GetTokenExpirationDate").Return(expDate).Once()
		mockedCfgSrv.On("GetRefreshTokenExpirationDate").Return(expDate).Once()
		authUser := domain.User{ID: 1}
		request.AddCookie(getRefreshTokenCookie(refreshToken))
		foundUser := domain.User{}
		mockedRepo.On("FindUserByID", authUser.ID).Return(&foundUser, nil).Once()

		recorder := httptest.NewRecorder()
		result := RefreshTokenHandler(recorder, request, h)

		okRes := results.CheckOkResult(t, result, http.StatusOK)
		resDto, isOk := okRes.Content.(*domain.TokenResponse)
		require.Equal(t, true, isOk, "should be a token response")
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTc0NzY0MDAsImlzQWRtaW4iOmZhbHNlLCJ1c2VySWQiOjAsInVzZXJOYW1lIjoiIn0.X2LZjUCGxdqgnUkpnXTkcZQuUSk7JgERVnQK4Vc6Sp0", resDto.Token)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTc0NzY0MDAsInVzZXJJZCI6MH0.jPcsnr6IZsQNPhpDB1--hW2EX1a1MCFrT0kujY6VQL4", resDto.RefreshToken)

		require.Equal(t, 1, len(recorder.Result().Cookies()))
		assert.Equal(t, "refreshToken", recorder.Result().Cookies()[0].Name)
		assert.Equal(t, resDto.RefreshToken, recorder.Result().Cookies()[0].Value)
		assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

		mockedRepo.AssertExpectations(t)
		mockedCfgSrv.AssertExpectations(t)
	})
}
