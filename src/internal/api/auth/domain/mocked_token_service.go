package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)

type MockedTokenService struct {
	mock.Mock
}

func NewMockedTokenService() *MockedTokenService {
	return &MockedTokenService{}
}

func (m *MockedTokenService) GenerateToken(user *UserEntity) (string, error) {
	args := m.Called(user)

	return args.String(0), args.Error(1)
}

func (m *MockedTokenService) GenerateRefreshToken(user *UserEntity, expirationDate time.Time) (string, error) {
	args := m.Called(user, expirationDate)

	return args.String(0), args.Error(1)
}

func (m *MockedTokenService) ParseToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockedTokenService) GetTokenInfo(token *jwt.Token) *TokenClaimsInfo {
	args := m.Called(token)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*TokenClaimsInfo)
}

func (m *MockedTokenService) GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo {
	args := m.Called(refreshToken)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*RefreshTokenClaimsInfo)
}
