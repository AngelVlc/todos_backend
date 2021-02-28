package helpers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
)

type TokenHelper interface {
	NewToken() interface{}
	GetTokenClaims(token interface{}) map[string]interface{}
	SignToken(token interface{}, secret string) (string, error)
	ParseToken(tokenString string, secret string) (interface{}, error)
	IsTokenValid(token interface{}) bool
}

type MockedTokenHelper struct {
	mock.Mock
}

func NewMockedTokenHelper() *MockedTokenHelper {
	return &MockedTokenHelper{}
}

func (m *MockedTokenHelper) NewToken() interface{} {
	args := m.Called()
	return args.Get(0).(interface{})
}

func (m *MockedTokenHelper) GetTokenClaims(token interface{}) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *MockedTokenHelper) SignToken(token interface{}, secret string) (string, error) {
	args := m.Called(token, secret)
	return args.String(0), args.Error(1)
}

func (m *MockedTokenHelper) ParseToken(tokenString string, secret string) (interface{}, error) {
	args := m.Called(tokenString, secret)

	got := args.Get(0)

	if got == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(interface{}), args.Error(1)
}

func (m *MockedTokenHelper) IsTokenValid(token interface{}) bool {
	args := m.Called(token)
	return args.Bool(0)
}

// JwtTokenHelper is the type used as JwtTokenHelper
type JwtTokenHelper struct {
	// secret string
}

// NewJwtTokenHelper returns a new JwtTokenHelper
func NewJwtTokenHelper() *JwtTokenHelper {
	return &JwtTokenHelper{}
}

// NewToken returns a new Jwt tooken
func (p *JwtTokenHelper) NewToken() interface{} {
	return jwt.New(jwt.SigningMethodHS256)
}

// GetTokenClaims returns the claims for the given token as a map
func (p *JwtTokenHelper) GetTokenClaims(token interface{}) map[string]interface{} {
	return p.getJwtToken(token).Claims.(jwt.MapClaims)
}

// SignToken signs the given token
func (p *JwtTokenHelper) SignToken(token interface{}, secret string) (string, error) {
	return p.getJwtToken(token).SignedString([]byte(secret))
}

// ParseToken parses the string and checks the signing method
func (p *JwtTokenHelper) ParseToken(tokenString string, secret string) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// IsTokenValid returns true if the given token is valid
func (p *JwtTokenHelper) IsTokenValid(token interface{}) bool {
	return p.getJwtToken(token).Valid
}

func (p *JwtTokenHelper) getJwtToken(token interface{}) *jwt.Token {
	jwtToken, _ := token.(*jwt.Token)
	return jwtToken
}
