package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	GenerateToken(user *User) (string, error)
	GenerateRefreshToken(user *User, expirationDate time.Time) (string, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetTokenInfo(token *jwt.Token) *TokenClaimsInfo
	GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo
}
