package services

import (
	"time"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
)

// AuthService is the service for auth methods
type AuthService struct {
	jwtPrv TokenHelper
	cfgSvc ConfigurationService
}

// NewAuthService returns a new auth service
func NewAuthService(jwtp TokenHelper, cfgSvc ConfigurationService) AuthService {
	return AuthService{jwtp, cfgSvc}
}

// GetTokens returns a new jwt token and a refresh token for the given user
func (s *AuthService) GetTokens(u *models.User) (map[string]string, error) {
	t := s.getNewToken(u)
	st, err := s.jwtPrv.SignToken(t, s.cfgSvc.GetJwtSecret())
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	rt := s.getNewRefreshToken(u)
	srt, err := s.jwtPrv.SignToken(rt, s.cfgSvc.GetJwtSecret())
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	result := map[string]string{
		"token":        st,
		"refreshToken": srt,
	}

	return result, nil
}

// ParseToken takes a token string, parses it and if it is valid returns a JwtClaimsInfo
// with its claims values
func (s *AuthService) ParseToken(tokenString string) (*models.JwtClaimsInfo, error) {
	token, err := s.jwtPrv.ParseToken(tokenString, s.cfgSvc.GetJwtSecret())
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid token", InternalError: err}
	}

	if !s.jwtPrv.IsTokenValid(token) {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid token"}
	}

	return s.getJwtInfo(token), nil
}

// ParseRefreshToken takes a refresh token string, parses it and if it is valid returns a
// RefreshTokenClaimsInfo with its claims values
func (s *AuthService) ParseRefreshToken(refreshTokenString string) (*models.RefreshTokenClaimsInfo, error) {
	refreshToken, err := s.jwtPrv.ParseToken(refreshTokenString, s.cfgSvc.GetJwtSecret())
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	if !s.jwtPrv.IsTokenValid(refreshToken) {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token"}
	}

	return s.getRefreshTokenInfo(refreshToken), nil
}

func (s *AuthService) getNewToken(u *models.User) interface{} {
	t := s.jwtPrv.NewToken()

	tc := s.jwtPrv.GetTokenClaims(t)
	tc["userName"] = u.Name
	tc["isAdmin"] = u.IsAdmin
	tc["userId"] = u.ID
	tc["exp"] = time.Now().Add(time.Minute * 15).Unix()

	return t
}

func (s *AuthService) getNewRefreshToken(u *models.User) interface{} {
	rt := s.jwtPrv.NewToken()
	rtc := s.jwtPrv.GetTokenClaims(rt)
	rtc["userId"] = u.ID
	rtc["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return rt
}

// GetJwtInfo returns a JwtClaimsInfo got from the token claims
func (s *AuthService) getJwtInfo(token interface{}) *models.JwtClaimsInfo {
	claims := s.jwtPrv.GetTokenClaims(token)

	info := models.JwtClaimsInfo{
		UserName: s.parseStringClaim(claims["userName"]),
		UserID:   s.parseInt32Claim(claims["userId"]),
		IsAdmin:  s.parseBoolClaim(claims["isAdmin"]),
	}

	return &info
}

/// GetRefreshTokenInfo returns a RefreshTokenClaimsInfo got from the refresh token claims
func (s *AuthService) getRefreshTokenInfo(refreshToken interface{}) *models.RefreshTokenClaimsInfo {
	claims := s.jwtPrv.GetTokenClaims(refreshToken)

	info := models.RefreshTokenClaimsInfo{
		UserID: s.parseInt32Claim(claims["userId"]),
	}
	return &info
}

func (s *AuthService) parseStringClaim(value interface{}) string {
	result, _ := value.(string)
	return result
}

func (s *AuthService) parseInt32Claim(value interface{}) int32 {
	result, _ := value.(float64)
	return int32(result)
}

func (s *AuthService) parseBoolClaim(value interface{}) bool {
	result, _ := value.(bool)
	return result
}
