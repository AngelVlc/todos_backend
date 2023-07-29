package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService interface {
	GenerateToken(user *UserRecord) (string, error)
	GenerateRefreshToken(user *UserRecord, expirationDate time.Time) (string, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetTokenInfo(token *jwt.Token) *TokenClaimsInfo
	GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo
}
