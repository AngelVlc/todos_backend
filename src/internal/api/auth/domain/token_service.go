package domain

import (
	"fmt"
	"time"

	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	cfgSvc sharedApp.ConfigurationService
}

func NewTokenService(cfgSvc sharedApp.ConfigurationService) *TokenService {
	return &TokenService{cfgSvc}
}

func (s *TokenService) GenerateToken(user *User) (string, error) {
	t := s.getNewToken(user.ID, user.Name, user.IsAdmin)
	return s.signToken(t, s.cfgSvc.GetJwtSecret())

}

func (s *TokenService) GenerateRefreshToken(user *User) (string, error) {
	rt := s.getNewRefreshToken(user.ID)
	return s.signToken(rt, s.cfgSvc.GetJwtSecret())
}

func (s *TokenService) getNewToken(userID int32, userName UserName, userIsAdmin bool) interface{} {
	t := s.newToken()

	tc := s.getTokenClaims(t)
	tc["userName"] = userName
	tc["isAdmin"] = userIsAdmin
	tc["userId"] = userID
	tc["exp"] = time.Now().Add(s.cfgSvc.TokenExpirationInSeconds()).Unix()

	return t
}

func (s *TokenService) getNewRefreshToken(userID int32) interface{} {
	rt := s.newToken()
	rtc := s.getTokenClaims(rt)
	rtc["userId"] = userID
	rtc["exp"] = time.Now().Add(s.cfgSvc.RefreshTokenExpirationInSeconds()).Unix()

	return rt
}

// newToken returns a new Jwt tooken
func (s *TokenService) newToken() interface{} {
	return jwt.New(jwt.SigningMethodHS256)
}

// getTokenClaims returns the claims for the given token as a map
func (s *TokenService) getTokenClaims(token interface{}) map[string]interface{} {
	return s.getJwtToken(token).Claims.(jwt.MapClaims)
}

// signToken signs the given token
func (s *TokenService) signToken(token interface{}, secret string) (string, error) {
	return s.getJwtToken(token).SignedString([]byte(secret))
}

// ParseToken parses the string and checks the signing method
func (s *TokenService) parseToken(tokenString string, secret string) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// IsTokenValid returns true if the given token is valid
func (s *TokenService) IsTokenValid(token interface{}) bool {
	return s.getJwtToken(token).Valid
}

func (s *TokenService) getJwtToken(token interface{}) *jwt.Token {
	jwtToken, _ := token.(*jwt.Token)
	return jwtToken
}

// ParseToken parses a refresh token string
func (s *TokenService) ParseToken(tokenString string) (interface{}, error) {
	return s.parseToken(tokenString, s.cfgSvc.GetJwtSecret())
}

// GetTokenInfo returns a JwtClaimsInfo got from the token claims
func (s *TokenService) GetTokenInfo(token interface{}) *JwtClaimsInfo {
	claims := s.getTokenClaims(token)

	info := JwtClaimsInfo{
		UserName: s.parseStringClaim(claims["userName"]),
		UserID:   s.parseInt32Claim(claims["userId"]),
		IsAdmin:  s.parseBoolClaim(claims["isAdmin"]),
	}

	return &info
}

/// GetRefreshTokenInfo returns a RefreshTokenClaimsInfo got from the refresh token claims
func (s *TokenService) GetRefreshTokenInfo(refreshToken interface{}) *RefreshTokenClaimsInfo {
	claims := s.getTokenClaims(refreshToken)

	info := RefreshTokenClaimsInfo{
		UserID: s.parseInt32Claim(claims["userId"]),
	}
	return &info
}

func (s *TokenService) parseStringClaim(value interface{}) string {
	result, _ := value.(string)
	return result
}

func (s *TokenService) parseInt32Claim(value interface{}) int32 {
	result, _ := value.(float64)
	return int32(result)
}

func (s *TokenService) parseBoolClaim(value interface{}) bool {
	result, _ := value.(bool)
	return result
}
