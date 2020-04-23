package providers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type TokenProvider interface {
	NewToken() *jwt.Token
	GetTokenClaims(token interface{}) map[string]interface{}
	SignToken(token interface{}, secret string) (string, error)
	ParseToken(tokenString string, secret string) (interface{}, error)
	IsTokenValid(token interface{}) bool
}

// JwtTokenProvider is the type used as JwtTokenProvider
type JwtTokenProvider struct {
	secret string
}

// NewJwtTokenProvider returns a new JwtTokenProvider
func NewJwtTokenProvider() *JwtTokenProvider {
	return &JwtTokenProvider{}
}

// NewToken returns a new Jwt tooken
func (p *JwtTokenProvider) NewToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

// GetTokenClaims returns the claims for the given token as a map
func (p *JwtTokenProvider) GetTokenClaims(token interface{}) map[string]interface{} {
	return p.getJwtToken(token).Claims.(jwt.MapClaims)
}

// SignToken signs the given token
func (p *JwtTokenProvider) SignToken(token interface{}, secret string) (string, error) {
	return p.getJwtToken(token).SignedString([]byte(secret))
}

// ParseToken parses the string and checks the signing method
func (p *JwtTokenProvider) ParseToken(tokenString string, secret string) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// IsTokenValid returns true if the given token is valid
func (p *JwtTokenProvider) IsTokenValid(token interface{}) bool {
	return p.getJwtToken(token).Valid
}

func (p *JwtTokenProvider) getJwtToken(token interface{}) *jwt.Token {
	jwtToken, _ := token.(*jwt.Token)
	return jwtToken
}
