package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AngelVlc/todos/internal/api/dtos"
	appErrors "github.com/AngelVlc/todos/internal/api/errors"
	"github.com/AngelVlc/todos/internal/api/models"
)

func TestAuthServiceGetTokens(t *testing.T) {
	jwtSecret := "jwtsecret"

	MockedTokenHelper := new(MockedTokenHelper)

	mockedCfgSvc := MockedConfigurationService{}
	mockedCfgSvc.On("GetJwtSecret").Return(jwtSecret)
	mockedCfgSvc.On("TokenExpirationInSeconds").Return(time.Second * 30)
	mockedCfgSvc.On("RefreshTokenExpirationInSeconds").Return(time.Hour)

	service := NewDefaultAuthService(MockedTokenHelper, &mockedCfgSvc)

	u := models.User{}
	token := struct{}{}
	claims := map[string]interface{}{}
	refreshToken := struct{}{}
	refreshTokenClaims := map[string]interface{}{}

	t.Run("should return an UnexpectedError if sign token fails", func(t *testing.T) {
		MockedTokenHelper.On("NewToken").Return(token).Once()
		MockedTokenHelper.On("GetTokenClaims", token).Return(claims).Once()
		MockedTokenHelper.On("SignToken", token, jwtSecret).Return("", fmt.Errorf("wadus")).Once()

		tokens, err := service.GetTokens(&u)

		assert.Nil(t, tokens)
		appErrors.CheckUnexpectedError(t, err, "Error creating jwt token", "wadus")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return an UnexpectedError if sign refresh token fails", func(t *testing.T) {
		MockedTokenHelper.On("NewToken").Return(token).Once()
		MockedTokenHelper.On("GetTokenClaims", token).Return(claims).Once()
		MockedTokenHelper.On("NewToken").Return(refreshToken).Once()
		MockedTokenHelper.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		MockedTokenHelper.On("SignToken", token, jwtSecret).Return("token", nil).Once()
		MockedTokenHelper.On("SignToken", refreshToken, jwtSecret).Return("", fmt.Errorf("wadus")).Once()

		tokens, err := service.GetTokens(&u)

		assert.Nil(t, tokens)
		appErrors.CheckUnexpectedError(t, err, "Error creating jwt refresh token", "wadus")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return the tokens if there are not errors", func(t *testing.T) {
		theToken := "theToken"
		theRefreshToken := "theRefreshToken"
		MockedTokenHelper.On("NewToken").Return(token).Once()
		MockedTokenHelper.On("GetTokenClaims", token).Return(claims).Once()
		MockedTokenHelper.On("NewToken").Return(refreshToken).Once()
		MockedTokenHelper.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		MockedTokenHelper.On("SignToken", token, jwtSecret).Return(theToken, nil).Once()
		MockedTokenHelper.On("SignToken", refreshToken, jwtSecret).Return(theRefreshToken, nil).Once()

		tokens, err := service.GetTokens(&u)

		want := &dtos.TokenResponseDto{
			Token:        theToken,
			RefreshToken: theRefreshToken,
		}

		assert.Equal(t, want, tokens)
		assert.Nil(t, err)

		MockedTokenHelper.AssertExpectations(t)
	})
}

func TestAuthServiceParseToken(t *testing.T) {
	jwtSecret := "jwtsecret"

	MockedTokenHelper := new(MockedTokenHelper)

	mockedCfgSvc := MockedConfigurationService{}
	mockedCfgSvc.On("GetJwtSecret").Return(jwtSecret)

	service := NewDefaultAuthService(MockedTokenHelper, &mockedCfgSvc)

	theToken := "theToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		MockedTokenHelper.On("ParseToken", theToken, jwtSecret).Return(nil, fmt.Errorf("wadus")).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid token", "wadus")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		token := struct{}{}

		MockedTokenHelper.On("ParseToken", theToken, jwtSecret).Return(token, nil).Once()
		MockedTokenHelper.On("IsTokenValid", token).Return(false).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid token", "")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return a jwt info when the token is valid", func(t *testing.T) {
		token := struct{}{}

		jwtInfo := models.JwtClaimsInfo{
			UserName: "wadus",
			IsAdmin:  true,
			UserID:   11,
		}

		MockedTokenHelper.On("ParseToken", theToken, jwtSecret).Return(token, nil).Once()
		MockedTokenHelper.On("IsTokenValid", token).Return(true).Once()

		c := map[string]interface{}{
			"userName": "wadus",
			"isAdmin":  true,
			"userId":   float64(11),
		}
		MockedTokenHelper.On("GetTokenClaims", token).Return(c).Once()

		res, err := service.ParseToken(theToken)

		assert.Equal(t, &jwtInfo, res)
		assert.Nil(t, err)
		MockedTokenHelper.AssertExpectations(t)
	})
}

func TestAuthServiceParseRefreshToken(t *testing.T) {
	jwtSecret := "jwtsecret"

	MockedTokenHelper := new(MockedTokenHelper)

	mockedCfgSvc := MockedConfigurationService{}
	mockedCfgSvc.On("GetJwtSecret").Return(jwtSecret)

	service := NewDefaultAuthService(MockedTokenHelper, &mockedCfgSvc)

	theRefreshToken := "theRefreshToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		MockedTokenHelper.On("ParseToken", theRefreshToken, jwtSecret).Return(nil, fmt.Errorf("wadus")).Once()

		jwtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid refresh token", "wadus")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		refreshToken := struct{}{}

		MockedTokenHelper.On("ParseToken", theRefreshToken, jwtSecret).Return(refreshToken, nil).Once()
		MockedTokenHelper.On("IsTokenValid", refreshToken).Return(false).Once()

		rtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, rtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid refresh token", "")
		MockedTokenHelper.AssertExpectations(t)
	})

	t.Run("should return a refresh token claims info when the token is valid", func(t *testing.T) {
		refreshToken := struct{}{}

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 11,
		}

		MockedTokenHelper.On("ParseToken", theRefreshToken, jwtSecret).Return(refreshToken, nil).Once()
		MockedTokenHelper.On("IsTokenValid", refreshToken).Return(true).Once()

		c := map[string]interface{}{
			"userId": float64(rtInfo.UserID),
		}
		MockedTokenHelper.On("GetTokenClaims", refreshToken).Return(c).Once()

		res, err := service.ParseRefreshToken(theRefreshToken)

		assert.Equal(t, &rtInfo, res)
		assert.Nil(t, err)
		MockedTokenHelper.AssertExpectations(t)
	})
}

func TestAuthServiceJwtProviderIntegration(t *testing.T) {
	jwtPrv := NewJwtTokenHelper()

	mockedEg := MockedEnvGetter{}
	mockedEg.On("Getenv", "JWT_SECRET").Return("jwtSecret")
	mockedEg.On("Getenv", "TOKEN_EXPIRATION_IN_SECONDS").Return("5m")
	mockedEg.On("Getenv", "REFRESH_TOKEN_EXPIRATION_IN_SECONDS").Return("1h")
	cfgSvc := NewDefaultConfigurationService(&mockedEg)

	service := NewDefaultAuthService(jwtPrv, cfgSvc)

	u := models.User{
		Name:    "wadus",
		IsAdmin: true,
		ID:      11,
	}

	tokens, err := service.GetTokens(&u)
	assert.NotNil(t, tokens)
	assert.Nil(t, err)

	jwtInfo, err := service.ParseToken(tokens.Token)
	assert.NotNil(t, jwtInfo)
	assert.Nil(t, err)

	assert.Equal(t, u.Name, jwtInfo.UserName)
	assert.Equal(t, u.IsAdmin, jwtInfo.IsAdmin)
	assert.Equal(t, u.ID, jwtInfo.UserID)

	rtClaims, err := service.ParseRefreshToken(tokens.RefreshToken)
	assert.NotNil(t, rtClaims)
	assert.Nil(t, err)

	assert.Equal(t, u.ID, rtClaims.UserID)
}
