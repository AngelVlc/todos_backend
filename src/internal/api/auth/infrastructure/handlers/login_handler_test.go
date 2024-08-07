//go:build !e2e
// +build !e2e

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_Returns_An_Error_If_The_Query_To_Find_The_User_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(nil, fmt.Errorf("some error")).Once()

	result := LoginHandler(httptest.NewRecorder(), request, h)

	results.CheckError(t, result, "some error")
	mockedAuthRepo.AssertExpectations(t)
}

func TestLoginHandler_Returns_An_ErrorResult_With_A_BadRequestError_If_The_Password_Does_Not_Match(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	foundUser := domain.UserRecord{PasswordHash: "hash"}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(&foundUser, nil).Once()

	result := LoginHandler(httptest.NewRecorder(), request, h)

	results.CheckBadRequestErrorResult(t, result, "Invalid password")
	mockedAuthRepo.AssertExpectations(t)
}

func TestLoginHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Generating_The_Token_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
	hashedPass := string(hashedBytes)
	foundUser := domain.UserRecord{ID: 1, PasswordHash: hashedPass}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(&foundUser, nil).Once()
	mockedTokenSrv.On("GenerateToken", foundUser.ToUserEntity()).Return("", fmt.Errorf("some error")).Once()

	result := LoginHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating jwt token")
	mockedAuthRepo.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestLoginHandler_Returns_An_ErrorResult_With_An_UnexpectedError_If_Generating_The_RefreshToken_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
	hashedPass := string(hashedBytes)
	foundUser := domain.UserRecord{ID: 1, PasswordHash: hashedPass}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(&foundUser, nil).Once()
	mockedTokenSrv.On("GenerateToken", foundUser.ToUserEntity()).Return("token", nil).Once()
	expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
	mockedCfgSrv.On("GetRefreshTokenExpirationTime").Return(expDate).Once()
	mockedTokenSrv.On("GenerateRefreshToken", foundUser.ToUserEntity(), expDate).Return("", fmt.Errorf("some error")).Once()

	result := LoginHandler(httptest.NewRecorder(), request, h)

	results.CheckUnexpectedErrorResult(t, result, "Error creating jwt refresh token")
	mockedAuthRepo.AssertExpectations(t)
	mockedCfgSrv.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestLoginHandler_Returns_An_OkResult_With_A_LoginResponse_And_Creates_The_Cookies_Although_Saving_The_RefreshToken_Fails(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
	hashedPass := string(hashedBytes)
	foundUser := domain.UserRecord{ID: 1, Name: "user", IsAdmin: true, PasswordHash: hashedPass}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(&foundUser, nil).Once()
	mockedTokenSrv.On("GenerateToken", foundUser.ToUserEntity()).Return("theToken", nil).Once()
	expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
	mockedCfgSrv.On("GetRefreshTokenExpirationTime").Return(expDate).Once()
	mockedTokenSrv.On("GenerateRefreshToken", foundUser.ToUserEntity(), expDate).Return("theRefreshToken", nil).Once()
	ctx := newrelic.NewContext(context.Background(), nil)
	mockedAuthRepo.On("CreateRefreshTokenIfNotExist", ctx, &domain.RefreshTokenEntity{UserID: foundUser.ID, RefreshToken: "theRefreshToken", ExpirationDate: expDate}).Return(fmt.Errorf("some error")).Once()

	recorder := httptest.NewRecorder()

	mockedAuthRepo.Wg.Add(1)
	result := LoginHandler(recorder, request, h)
	mockedAuthRepo.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a user response")
	assert.Equal(t, int32(1), res.ID)
	assert.Equal(t, "user", res.Name)
	assert.True(t, res.IsAdmin)

	require.Equal(t, 2, len(recorder.Result().Cookies()))
	assert.Equal(t, "token", recorder.Result().Cookies()[0].Name)
	assert.Equal(t, "theToken", recorder.Result().Cookies()[0].Value)
	assert.Equal(t, "refreshToken", recorder.Result().Cookies()[1].Name)
	assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[1].Value)
	assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

	mockedAuthRepo.AssertExpectations(t)
	mockedCfgSrv.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}

func TestLoginHandler_Returns_An_OkResult_With_A_LoginResponse_And_Creates_The_Cookies_And_Saves_The_RefreshToken(t *testing.T) {
	mockedAuthRepo := repository.MockedAuthRepository{}
	mockedUsersRepo := repository.MockedUsersRepository{}
	mockedCfgSrv := application.MockedConfigurationService{}
	mockedTokenSrv := domain.MockedTokenService{}
	userName, _ := domain.NewUserNameValueObject("wadus")
	userPassword, _ := domain.NewUserPasswordValueObject("pass")
	h := handler.Handler{
		AuthRepository:  &mockedAuthRepo,
		UsersRepository: &mockedUsersRepo,
		CfgSrv:          &mockedCfgSrv,
		TokenSrv:        &mockedTokenSrv,
		RequestInput:    &infrastructure.LoginInput{UserName: userName, Password: userPassword},
	}

	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("pass"), 10)
	hashedPass := string(hashedBytes)
	foundUser := domain.UserRecord{ID: 1, Name: "user", IsAdmin: true, PasswordHash: hashedPass}
	mockedUsersRepo.On("FindUser", request.Context(), domain.UserRecord{Name: "wadus"}).Return(&foundUser, nil).Once()
	mockedTokenSrv.On("GenerateToken", foundUser.ToUserEntity()).Return("theToken", nil).Once()
	expDate, _ := time.Parse(time.RFC3339, "2021-04-03T19:00:00+00:00")
	mockedCfgSrv.On("GetRefreshTokenExpirationTime").Return(expDate).Once()
	mockedTokenSrv.On("GenerateRefreshToken", foundUser.ToUserEntity(), expDate).Return("theRefreshToken", nil).Once()
	ctx := newrelic.NewContext(context.Background(), nil)
	mockedAuthRepo.On("CreateRefreshTokenIfNotExist", ctx, &domain.RefreshTokenEntity{UserID: foundUser.ID, RefreshToken: "theRefreshToken", ExpirationDate: expDate}).Return(nil).Once()

	recorder := httptest.NewRecorder()

	mockedAuthRepo.Wg.Add(1)
	result := LoginHandler(recorder, request, h)
	mockedAuthRepo.Wg.Wait()

	okRes := results.CheckOkResult(t, result, http.StatusOK)
	res, isOk := okRes.Content.(infrastructure.UserResponse)
	require.Equal(t, true, isOk, "should be a login response")
	assert.Equal(t, int32(1), res.ID)
	assert.Equal(t, "user", res.Name)
	assert.True(t, res.IsAdmin)

	require.Equal(t, 2, len(recorder.Result().Cookies()))
	assert.Equal(t, "token", recorder.Result().Cookies()[0].Name)
	assert.Equal(t, "theToken", recorder.Result().Cookies()[0].Value)
	assert.Equal(t, "refreshToken", recorder.Result().Cookies()[1].Name)
	assert.Equal(t, "theRefreshToken", recorder.Result().Cookies()[1].Value)
	assert.True(t, recorder.Result().Cookies()[0].HttpOnly)

	mockedAuthRepo.AssertExpectations(t)
	mockedCfgSrv.AssertExpectations(t)
	mockedTokenSrv.AssertExpectations(t)
}
