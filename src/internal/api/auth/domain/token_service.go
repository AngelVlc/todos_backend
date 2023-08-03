package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService interface {
	GenerateToken(user *UserEntity) (string, error)
	GenerateRefreshToken(user *UserEntity, expirationDate time.Time) (string, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetTokenInfo(token *jwt.Token) *TokenClaimsInfo
	GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo
}
