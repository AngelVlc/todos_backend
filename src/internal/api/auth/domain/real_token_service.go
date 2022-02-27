package domain

import (
	"fmt"
	"time"

	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/golang-jwt/jwt"
)

type RealTokenService struct {
	cfgSvc sharedApp.ConfigurationService
}

func NewRealTokenService(cfgSvc sharedApp.ConfigurationService) *RealTokenService {
	return &RealTokenService{cfgSvc}
}

func (s *RealTokenService) GenerateToken(user *User) (string, error) {
	t := s.getNewToken(user.ID, user.Name, user.IsAdmin)

	return s.signToken(t, s.cfgSvc.GetJwtSecret())

}

func (s *RealTokenService) GenerateRefreshToken(user *User, expirationDate time.Time) (string, error) {
	rt := s.getNewRefreshToken(user.ID, expirationDate)

	return s.signToken(rt, s.cfgSvc.GetJwtSecret())
}

// ParseToken parses a refresh token string
func (s *RealTokenService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := s.cfgSvc.GetJwtSecret()

		return []byte(secret), nil
	})
}

// GetTokenInfo returns a JwtClaimsInfo got from the token claims
func (s *RealTokenService) GetTokenInfo(token *jwt.Token) *TokenClaimsInfo {
	claims := s.getTokenClaims(token)

	info := TokenClaimsInfo{
		UserName: s.parseStringClaim(claims["userName"]),
		UserID:   s.parseInt32Claim(claims["userId"]),
		IsAdmin:  s.parseBoolClaim(claims["isAdmin"]),
	}

	return &info
}

/// GetRefreshTokenInfo returns a RefreshTokenClaimsInfo got from the refresh token claims
func (s *RealTokenService) GetRefreshTokenInfo(refreshToken *jwt.Token) *RefreshTokenClaimsInfo {
	claims := s.getTokenClaims(refreshToken)

	info := RefreshTokenClaimsInfo{
		UserID: s.parseInt32Claim(claims["userId"]),
	}

	return &info
}

func (s *RealTokenService) getNewToken(userID int32, userName UserName, userIsAdmin bool) *jwt.Token {
	t := s.newToken()

	tc := s.getTokenClaims(t)
	tc["userName"] = userName
	tc["isAdmin"] = userIsAdmin
	tc["userId"] = userID
	tc["exp"] = s.cfgSvc.GetTokenExpirationDate().Unix()

	return t
}

func (s *RealTokenService) getNewRefreshToken(userID int32, expirationDate time.Time) *jwt.Token {
	rt := s.newToken()
	rtc := s.getTokenClaims(rt)
	rtc["userId"] = userID
	rtc["exp"] = expirationDate.Unix()

	return rt
}

// newToken returns a new Jwt tooken
func (s *RealTokenService) newToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

// getTokenClaims returns the claims for the given token as a map
func (s *RealTokenService) getTokenClaims(token *jwt.Token) map[string]interface{} {
	return token.Claims.(jwt.MapClaims)
}

// signToken signs the given token
func (s *RealTokenService) signToken(token *jwt.Token, secret string) (string, error) {
	return token.SignedString([]byte(secret))
}

func (s *RealTokenService) parseStringClaim(value interface{}) string {
	result, _ := value.(string)

	return result
}

func (s *RealTokenService) parseInt32Claim(value interface{}) int32 {
	result, _ := value.(float64)

	return int32(result)
}

func (s *RealTokenService) parseBoolClaim(value interface{}) bool {
	result, _ := value.(bool)

	return result
}
