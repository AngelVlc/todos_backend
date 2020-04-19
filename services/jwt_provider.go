package services

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JwtProvider is the type used as JwtProvider
type JwtProvider struct {
	secret string
}

// NewJwtProvider returns a new JwtProvider
func NewJwtProvider(cfgSvc ConfigurationService) JwtProvider {
	return JwtProvider{cfgSvc.GetJwtSecret()}
}

// NewToken returns a new Jwt tooken
func (p *JwtProvider) NewToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

// GetTokenClaims returns the claims for the given token as a map
func (p *JwtProvider) GetTokenClaims(token interface{}) map[string]interface{} {
	return p.getJwtToken(token).Claims.(jwt.MapClaims)
}

// SignToken signs the given token
func (p *JwtProvider) SignToken(token interface{}) (string, error) {
	return p.getJwtToken(token).SignedString([]byte(p.secret))
}

// ParseToken parses the string and checks the signing method
func (p *JwtProvider) ParseToken(tokenString string) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
}

// IsTokenValid returns true if the given token is valid
func (p *JwtProvider) IsTokenValid(token interface{}) bool {
	return p.getJwtToken(token).Valid
}

func (p *JwtProvider) getJwtToken(token interface{}) *jwt.Token {
	jwtToken, _ := token.(*jwt.Token)
	return jwtToken
}
