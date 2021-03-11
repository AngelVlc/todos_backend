//+build !e2e

package infrastructure

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	t.Run("Should return an errorResult with an UnauthorizedError if the refresh token is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		request.AddCookie(getRefreshTokenCookie("badToken"))
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "Error parsing the refresh token")
		mockedCfgSrv.AssertExpectations(t)

	})

	t.Run("Should return an errorResult with an UnexpectedError if getting the user by id fails", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(2)
		mockedCfgSrv.On("RefreshTokenExpirationInSeconds").Return(5 * time.Minute).Once()
		authUser := domain.User{ID: 1}
		rt, _ := domain.NewTokenService(&mockedCfgSrv).GenerateRefreshToken(&authUser)
		request.AddCookie(getRefreshTokenCookie(rt))
		mockedRepo.On("FindUserByID", &authUser.ID).Return(nil, fmt.Errorf("some error")).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnexpectedErrorResult(t, result, "Error getting user by user id")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an errorResult with an UnauthorizedError if the user no longer exists", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(2)
		mockedCfgSrv.On("RefreshTokenExpirationInSeconds").Return(5 * time.Minute).Once()
		authUser := domain.User{ID: 1}
		rt, _ := domain.NewTokenService(&mockedCfgSrv).GenerateRefreshToken(&authUser)
		request.AddCookie(getRefreshTokenCookie(rt))
		mockedRepo.On("FindUserByID", &authUser.ID).Return(nil, nil).Once()

		result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

		results.CheckUnauthorizedErrorErrorResult(t, result, "The user no longer exists")
		mockedCfgSrv.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Should return an okResult with the tokens and should create the cookie if the refresh token is correct", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(4)
		mockedCfgSrv.On("TokenExpirationInSeconds").Return(5 * time.Minute).Once()
		mockedCfgSrv.On("RefreshTokenExpirationInSeconds").Return(5 * time.Minute).Times(2)
		authUser := domain.User{ID: 1}
		rt, _ := domain.NewTokenService(&mockedCfgSrv).GenerateRefreshToken(&authUser)
		request.AddCookie(getRefreshTokenCookie(rt))
		foundUser := domain.User{}
		mockedRepo.On("FindUserByID", &authUser.ID).Return(&foundUser, nil).Once()

		recorder := httptest.NewRecorder()
		result := RefreshTokenHandler(recorder, request, h)

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
