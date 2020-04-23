package services

// import (
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/AngelVlc/lists-backend/models"
// 	appErrors "github.com/AngelVlc/todos/errors"
// 	"github.com/stretchr/testify/mock"
// )

// type mockedJwtProvider struct {
// 	mock.Mock
// }

// func (m *mockedJwtProvider) NewToken() interface{} {
// 	args := m.Called()
// 	return args.Get(0).(interface{})
// }

// func (m *mockedJwtProvider) GetTokenClaims(token interface{}) map[string]interface{} {
// 	args := m.Called(token)
// 	return args.Get(0).(map[string]interface{})
// }

// func (m *mockedJwtProvider) SignToken(token interface{}, secret string) (string, error) {
// 	args := m.Called(token, secret)
// 	return args.String(0), args.Error(1)
// }

// func (m *mockedJwtProvider) ParseToken(tokenString string, secret string) (interface{}, error) {
// 	args := m.Called(tokenString, secret)

// 	got := args.Get(0)

// 	if got == nil {
// 		return nil, args.Error(1)
// 	}

// 	return args.Get(0).(interface{}), args.Error(1)
// }

// func (m *mockedJwtProvider) IsTokenValid(token interface{}) bool {
// 	args := m.Called(token)
// 	return args.Bool(0)
// }

// func TestAuthServiceGetTokens(t *testing.T) {
// 	mockedJwtProvider := new(mockedJwtProvider)

// 	mockedEg := MockedEnvGetter{}
// 	jwtSecret := "jwtsecret"
// 	mockedEg.On("Getenv", "JWT_SECRET").Return(jwtSecret)
// 	cfgSvc := NewConfigurationService(&mockedEg)

// 	service := NewAuthService(mockedJwtProvider, cfgSvc)

// 	u := models.User{}
// 	token := struct{}{}
// 	claims := map[string]interface{}{}
// 	refreshToken := struct{}{}
// 	refreshTokenClaims := map[string]interface{}{}

// 	t.Run("should return an UnexpectedError if sign token fails", func(t *testing.T) {
// 		mockedJwtProvider.On("NewToken").Return(token).Once()
// 		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
// 		mockedJwtProvider.On("SignToken", token, jwtSecret).Return("", errors.New("wadus")).Once()

// 		tokens, err := service.GetTokens(&u)

// 		assert.Nil(t, tokens)
// 		assert.NotNil(t, err)
// 		unexpectedErr, isUnexpectedErr := err.(*appErrors.UnexpectedError)
// 		assert.Equal(t, true, isUnexpectedErr, "should be an unexpected error")
// 		assert.Equal(t, "Error creating jwt token", unexpectedErr.Error())
// 		mockedJwtProvider.AssertExpectations(t)
// 	})

// 	// t.Run("should return an UnexpectedError if sign refresh token fails", func(t *testing.T) {
// 	// 	mockedJwtProvider.On("NewToken").Return(token).Once()
// 	// 	mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
// 	// 	mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
// 	// 	mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
// 	// 	mockedJwtProvider.On("SignToken", token).Return("token", nil).Once()
// 	// 	mockedJwtProvider.On("SignToken", refreshToken).Return("", errors.New("wadus")).Once()

// 	// 	tokens, err := service.CreateTokens(&u)

// 	// 	assert.Nil(t, tokens)
// 	// 	assert.NotNil(t, err)
// 	// 	unexpectedErr, isUnexpectedErr := err.(*appErrors.UnexpectedError)
// 	// 	assert.Equal(t, true, isUnexpectedErr, "should be an unexpected error")
// 	// 	assert.Equal(t, "Error creating jwt refresh token", unexpectedErr.Error())
// 	// 	mockedJwtProvider.AssertExpectations(t)
// 	// })

// 	// t.Run("should return a signed token if no error happen", func(t *testing.T) {
// 	// 	theToken := "theToken"
// 	// 	theRefreshToken := "theRefreshToken"
// 	// 	mockedJwtProvider.On("NewToken").Return(token).Once()
// 	// 	mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
// 	// 	mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
// 	// 	mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
// 	// 	mockedJwtProvider.On("SignToken", token).Return(theToken, nil).Once()
// 	// 	mockedJwtProvider.On("SignToken", refreshToken).Return(theRefreshToken, nil).Once()

// 	// 	tokens, err := service.CreateTokens(&u)

// 	// 	want := map[string]string{
// 	// 		"token":        theToken,
// 	// 		"refreshToken": theRefreshToken,
// 	// 	}

// 	// 	assert.Equal(t, want, tokens)
// 	// 	assert.Nil(t, err)

// 	// 	mockedJwtProvider.AssertExpectations(t)
// 	// })
// }
