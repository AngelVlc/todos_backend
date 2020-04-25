package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/stretchr/testify/mock"
)

type mockedJwtProvider struct {
	mock.Mock
}

func (m *mockedJwtProvider) NewToken() interface{} {
	args := m.Called()
	return args.Get(0).(interface{})
}

func (m *mockedJwtProvider) GetTokenClaims(token interface{}) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *mockedJwtProvider) SignToken(token interface{}, secret string) (string, error) {
	args := m.Called(token, secret)
	return args.String(0), args.Error(1)
}

func (m *mockedJwtProvider) ParseToken(tokenString string, secret string) (interface{}, error) {
	args := m.Called(tokenString, secret)

	got := args.Get(0)

	if got == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(interface{}), args.Error(1)
}

func (m *mockedJwtProvider) IsTokenValid(token interface{}) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func TestAuthServiceGetTokens(t *testing.T) {
	jwtSecret := "jwtsecret"

	mockedJwtProvider := new(mockedJwtProvider)

	mockedEg := MockedEnvGetter{}
	mockedEg.On("Getenv", "JWT_SECRET").Return(jwtSecret)
	cfgSvc := NewConfigurationService(&mockedEg)

	service := NewDefaultAuthService(mockedJwtProvider, cfgSvc)

	u := models.User{}
	token := struct{}{}
	claims := map[string]interface{}{}
	refreshToken := struct{}{}
	refreshTokenClaims := map[string]interface{}{}

	t.Run("should return an UnexpectedError if sign token fails", func(t *testing.T) {
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("SignToken", token, jwtSecret).Return("", fmt.Errorf("wadus")).Once()

		tokens, err := service.GetTokens(&u)

		assert.Nil(t, tokens)
		appErrors.CheckUnexpectedError(t, err, "Error creating jwt token", "wadus")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an UnexpectedError if sign refresh token fails", func(t *testing.T) {
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		mockedJwtProvider.On("SignToken", token, jwtSecret).Return("token", nil).Once()
		mockedJwtProvider.On("SignToken", refreshToken, jwtSecret).Return("", fmt.Errorf("wadus")).Once()

		tokens, err := service.GetTokens(&u)

		assert.Nil(t, tokens)
		appErrors.CheckUnexpectedError(t, err, "Error creating jwt refresh token", "wadus")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return the tokens if there are not errors", func(t *testing.T) {
		theToken := "theToken"
		theRefreshToken := "theRefreshToken"
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		mockedJwtProvider.On("SignToken", token, jwtSecret).Return(theToken, nil).Once()
		mockedJwtProvider.On("SignToken", refreshToken, jwtSecret).Return(theRefreshToken, nil).Once()

		tokens, err := service.GetTokens(&u)

		want := map[string]string{
			"token":        theToken,
			"refreshToken": theRefreshToken,
		}

		assert.Equal(t, want, tokens)
		assert.Nil(t, err)

		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceParseToken(t *testing.T) {
	jwtSecret := "jwtsecret"

	mockedJwtProvider := new(mockedJwtProvider)

	mockedEg := MockedEnvGetter{}
	mockedEg.On("Getenv", "JWT_SECRET").Return(jwtSecret)
	cfgSvc := NewConfigurationService(&mockedEg)

	service := NewDefaultAuthService(mockedJwtProvider, cfgSvc)

	theToken := "theToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		mockedJwtProvider.On("ParseToken", theToken, jwtSecret).Return(nil, fmt.Errorf("wadus")).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid token", "wadus")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		token := struct{}{}

		mockedJwtProvider.On("ParseToken", theToken, jwtSecret).Return(token, nil).Once()
		mockedJwtProvider.On("IsTokenValid", token).Return(false).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid token", "")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return a jwt info when the token is valid", func(t *testing.T) {
		token := struct{}{}

		jwtInfo := models.JwtClaimsInfo{
			UserName: "wadus",
			IsAdmin:  true,
			UserID:   11,
		}

		mockedJwtProvider.On("ParseToken", theToken, jwtSecret).Return(token, nil).Once()
		mockedJwtProvider.On("IsTokenValid", token).Return(true).Once()

		c := map[string]interface{}{
			"userName": "wadus",
			"isAdmin":  true,
			"userId":   float64(11),
		}
		mockedJwtProvider.On("GetTokenClaims", token).Return(c).Once()

		res, err := service.ParseToken(theToken)

		assert.Equal(t, &jwtInfo, res)
		assert.Nil(t, err)
		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceParseRefreshToken(t *testing.T) {
	jwtSecret := "jwtsecret"

	mockedJwtProvider := new(mockedJwtProvider)

	mockedEg := MockedEnvGetter{}
	mockedEg.On("Getenv", "JWT_SECRET").Return(jwtSecret)
	cfgSvc := NewConfigurationService(&mockedEg)

	service := NewDefaultAuthService(mockedJwtProvider, cfgSvc)

	theRefreshToken := "theRefreshToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		mockedJwtProvider.On("ParseToken", theRefreshToken, jwtSecret).Return(nil, fmt.Errorf("wadus")).Once()

		jwtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, jwtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid refresh token", "wadus")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		refreshToken := struct{}{}

		mockedJwtProvider.On("ParseToken", theRefreshToken, jwtSecret).Return(refreshToken, nil).Once()
		mockedJwtProvider.On("IsTokenValid", refreshToken).Return(false).Once()

		rtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, rtInfo)
		appErrors.CheckUnathorizedError(t, err, "Invalid refresh token", "")
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return a refresh token claims info when the token is valid", func(t *testing.T) {
		refreshToken := struct{}{}

		rtInfo := models.RefreshTokenClaimsInfo{
			UserID: 11,
		}

		mockedJwtProvider.On("ParseToken", theRefreshToken, jwtSecret).Return(refreshToken, nil).Once()
		mockedJwtProvider.On("IsTokenValid", refreshToken).Return(true).Once()

		c := map[string]interface{}{
			"userId": float64(rtInfo.UserID),
		}
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(c).Once()

		res, err := service.ParseRefreshToken(theRefreshToken)

		assert.Equal(t, &rtInfo, res)
		assert.Nil(t, err)
		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceJwtProviderIntegration(t *testing.T) {
	jwtPrv := NewJwtTokenHelper()

	mockedEg := MockedEnvGetter{}
	mockedEg.On("Getenv", "JWT_SECRET").Return("jwtSecret")
	cfgSvc := NewConfigurationService(&mockedEg)

	service := NewDefaultAuthService(jwtPrv, cfgSvc)

	u := models.User{
		Name:    "wadus",
		IsAdmin: true,
		ID:      11,
	}

	tokens, err := service.GetTokens(&u)
	assert.NotNil(t, tokens)
	assert.Nil(t, err)

	jwtInfo, err := service.ParseToken(tokens["token"])
	assert.NotNil(t, jwtInfo)
	assert.Nil(t, err)

	assert.Equal(t, u.Name, jwtInfo.UserName)
	assert.Equal(t, u.IsAdmin, jwtInfo.IsAdmin)
	assert.Equal(t, u.ID, jwtInfo.UserID)

	rtClaims, err := service.ParseRefreshToken(tokens["refreshToken"])
	assert.NotNil(t, rtClaims)
	assert.Nil(t, err)

	assert.Equal(t, u.ID, rtClaims.UserID)
}
