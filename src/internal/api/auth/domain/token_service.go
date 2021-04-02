package domain

import (
	"fmt"

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

// ParseToken parses a refresh token string
func (s *TokenService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		secret := s.cfgSvc.GetJwtSecret()
		return []byte(secret), nil
	})
}

// GetTokenInfo returns a JwtClaimsInfo got from the token claims
func (s *TokenService) GetTokenInfo(token *jwt.Token) *TokenClaimsInfo {
	claims := s.getTokenClaims(token)

	info := TokenClaimsInfo{
		UserName: s.parseStringClaim(claims["userName"]),
		UserID:   s.parseInt32Claim(claims["userId"]),
		IsAdmin:  s.parseBoolClaim(claims["isAdmin"]),
	}

	return &info
}

/// GetRefreshTokenInfo returns a RefreshTokenClaimsInfo got from the refresh token claims
func (s *TokenService) GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo {
	claims := s.getTokenClaims(refreshToken)

	info := RefreshTokenClaimsInfo{
		UserID: s.parseInt32Claim(claims["userId"]),
	}
	return &info
}

func (s *TokenService) getNewToken(userID int32, userName UserName, userIsAdmin bool) *jwt.Token {
	t := s.newToken()

	tc := s.getTokenClaims(t)
	tc["userName"] = userName
	tc["isAdmin"] = userIsAdmin
	tc["userId"] = userID
	tc["exp"] = s.cfgSvc.GetTokenExpirationDate().Unix()

	return t
}

func (s *TokenService) getNewRefreshToken(userID int32) *jwt.Token {
	rt := s.newToken()
	rtc := s.getTokenClaims(rt)
	rtc["userId"] = userID
	rtc["exp"] = s.cfgSvc.GetRefreshTokenExpirationDate().Unix()

	return rt
}

// newToken returns a new Jwt tooken
func (s *TokenService) newToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

// getTokenClaims returns the claims for the given token as a map
func (s *TokenService) getTokenClaims(token *jwt.Token) map[string]interface{} {
	return token.Claims.(jwt.MapClaims)
}

// signToken signs the given token
func (s *TokenService) signToken(token *jwt.Token, secret string) (string, error) {
	return token.SignedString([]byte(secret))
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
