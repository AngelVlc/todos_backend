//go:build !e2e
// +build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenHandler_Validations_Returns_An_ErrorResult_With_A_BadRequestError_If_There_Is_Not_A_Refresh_Token_Cookie(t *testing.T) {
	h := handler.Handler{}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Missing refresh token cookie")
}

func TestRefreshTokenHandler_Returns_An_ErrorResult_With_An_UnauthorizedError_If_The_RefreshToken_Is_Not_Valid(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	mockedTokenSrv.On("ParseToken", "badToken").Return(nil, fmt.Errorf("some error")).Once()

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("badToken"))

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckUnauthorizedErrorErrorResult(t, result, "Invalid refresh token")
	mockedTokenSrv.AssertExpectations(t)
}

func TestRefreshTokenHandler_Returns_An_Error_If_Getting_The_User_By_Id_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	token := jwt.Token{Valid: true}
	mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
	rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
	mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("token"))
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserEntity{ID: 1}).Return(nil, fmt.Errorf("some error")).Once()

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedCfgSrv.AssertExpectations(t)
	mockedAuthRepo.AssertExpectations(t)
	mockedUsersRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestRefreshTokenHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Checking_If_The_RefreshToken_Exists_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	token := jwt.Token{Valid: true}
	mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
	rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
	mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("token"))
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserEntity{ID: 1}).Return(&domain.UserEntity{}, nil).Once()
	mockedAuthRepo.On("ExistsRefreshToken", request.Context(), domain.RefreshTokenEntity{RefreshToken: "token", UserID: 1}).Return(false, fmt.Errorf("some error")).Once()

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error getting the refresh token")
	mockedCfgSrv.AssertExpectations(t)
	mockedAuthRepo.AssertExpectations(t)
	mockedUsersRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestRefreshTokenHandler_Returns_An_ErrorResult_With_An_UnauthorizedError_If_The_RefreshToken_Does_Not_Exist(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	token := jwt.Token{Valid: true}
	mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
	rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
	mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("token"))
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserEntity{ID: 1}).Return(&domain.UserEntity{}, nil).Once()
	mockedAuthRepo.On("ExistsRefreshToken", request.Context(), domain.RefreshTokenEntity{RefreshToken: "token", UserID: 1}).Return(false, nil).Once()

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckUnauthorizedErrorErrorResult(t, result, "The refresh token is not valid")
	mockedCfgSrv.AssertExpectations(t)
	mockedAuthRepo.AssertExpectations(t)
	mockedUsersRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestRefreshTokenHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Generate_The_New_Token_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	token := jwt.Token{Valid: true}
	mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
	rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
	mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("token"))
	foundUser := domain.UserEntity{}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserEntity{ID: 1}).Return(&foundUser, nil).Once()
	mockedAuthRepo.On("ExistsRefreshToken", request.Context(), domain.RefreshTokenEntity{RefreshToken: "token", UserID: 1}).Return(true, nil).Once()
	mockedTokenSrv.On("GenerateToken", &foundUser).Return("", fmt.Errorf("some error")).Once()

	result := RefreshTokenHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating jwt token")
	mockedCfgSrv.AssertExpectations(t)
	mockedAuthRepo.AssertExpectations(t)
	mockedUsersRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestRefreshTokenHandler_Returns_An_OkResult_With_The_Token_And_Creates_The_Cookie(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	h := handler.Handler{AuthRepository: &mockedAuthRepo, UsersRepository: &mockedUsersRepo, CfgSrv: &mockedCfgSrv, TokenSrv: &mockedTokenSrv}

	getRefreshTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: refreshTokenCookieName, Value: rt}
	}

	token := jwt.Token{Valid: true}
	mockedTokenSrv.On("ParseToken", "token").Return(&token, nil).Once()
	rtClaims := domain.RefreshTokenClaimsInfo{UserID: 1}
	mockedTokenSrv.On("GetRefreshTokenInfo", &token).Return(&rtClaims).Once()
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(getRefreshTokenCookie("token"))
	foundUser := domain.UserEntity{}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserEntity{ID: 1}).Return(&foundUser, nil).Once()
	mockedAuthRepo.On("ExistsRefreshToken", request.Context(), domain.RefreshTokenEntity{RefreshToken: "token", UserID: 1}).Return(true, nil).Once()
	mockedTokenSrv.On("GenerateToken", &foundUser).Return("theToken", nil).Once()

	recorder := httptest.NewRecorder()
	result := RefreshTokenHandler(recorder, request, h)

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	assert.Nil(t, okRes.Content)

	require.Equal(t, 1, len(recorder.Result().Cookies()))
	assert.Equal(t, "token", recorder.Result().Cookies()[0].Name)
	assert.Equal(t, "theToken", recorder.Result().Cookies()[0].Value)
	assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

	mockedAuthRepo.AssertExpectations(t)
	mockedCfgSrv.AssertExpectations(t)
	mockedUsersRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}
