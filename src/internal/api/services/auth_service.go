package services

import (
	"time"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/models"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/stretchr/testify/mock"
)

type AuthService interface {
	GetTokens(u *models.User) (*dtos.TokenResponseDto, error)
	ParseToken(tokenString string) (*models.JwtClaimsInfo, error)
	ParseRefreshToken(refreshTokenString string) (*models.RefreshTokenClaimsInfo, error)
}

type MockedAuthService struct {
	mock.Mock
}

func NewMockedAuthService() *MockedAuthService {
	return &MockedAuthService{}
}

func (m *MockedAuthService) GetTokens(u *models.User) (*dtos.TokenResponseDto, error) {
	args := m.Called(u)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.TokenResponseDto), args.Error(1)
}

func (m *MockedAuthService) ParseToken(t string) (*models.JwtClaimsInfo, error) {
	args := m.Called(t)

	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JwtClaimsInfo), args.Error(1)
}

func (m *MockedAuthService) ParseRefreshToken(t string) (*models.RefreshTokenClaimsInfo, error) {
	args := m.Called(t)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshTokenClaimsInfo), args.Error(1)
}

// DefaultAuthService is the service for auth methods
type DefaultAuthService struct {
	jwtPrv TokenHelper
	cfgSvc ConfigurationService
}

// NewDefaultAuthService returns a new auth service
func NewDefaultAuthService(jwtp TokenHelper, cfgSvc ConfigurationService) *DefaultAuthService {
	return &DefaultAuthService{jwtp, cfgSvc}
}

// GetTokens returns a new jwt token and a refresh token for the given user
func (s *DefaultAuthService) GetTokens(u *models.User) (*dtos.TokenResponseDto, error) {
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

	result := dtos.TokenResponseDto{Token: st, RefreshToken: srt}

	return &result, nil
}

// ParseToken takes a token string, parses it and if it is valid returns a JwtClaimsInfo
// with its claims values
func (s *DefaultAuthService) ParseToken(tokenString string) (*models.JwtClaimsInfo, error) {
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
func (s *DefaultAuthService) ParseRefreshToken(refreshTokenString string) (*models.RefreshTokenClaimsInfo, error) {
	refreshToken, err := s.jwtPrv.ParseToken(refreshTokenString, s.cfgSvc.GetJwtSecret())
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	if !s.jwtPrv.IsTokenValid(refreshToken) {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token"}
	}

	return s.getRefreshTokenInfo(refreshToken), nil
}

func (s *DefaultAuthService) getNewToken(u *models.User) interface{} {
	t := s.jwtPrv.NewToken()

	tc := s.jwtPrv.GetTokenClaims(t)
	tc["userName"] = u.Name
	tc["isAdmin"] = u.IsAdmin
	tc["userId"] = u.ID
	tc["exp"] = time.Now().Add(s.cfgSvc.TokenExpirationInSeconds()).Unix()

	return t
}

func (s *DefaultAuthService) getNewRefreshToken(u *models.User) interface{} {
	rt := s.jwtPrv.NewToken()
	rtc := s.jwtPrv.GetTokenClaims(rt)
	rtc["userId"] = u.ID
	rtc["exp"] = time.Now().Add(s.cfgSvc.RefreshTokenExpirationInSeconds()).Unix()

	return rt
}

// GetJwtInfo returns a JwtClaimsInfo got from the token claims
func (s *DefaultAuthService) getJwtInfo(token interface{}) *models.JwtClaimsInfo {
	claims := s.jwtPrv.GetTokenClaims(token)

	info := models.JwtClaimsInfo{
		UserName: s.parseStringClaim(claims["userName"]),
		UserID:   s.parseInt32Claim(claims["userId"]),
		IsAdmin:  s.parseBoolClaim(claims["isAdmin"]),
	}

	return &info
}

/// GetRefreshTokenInfo returns a RefreshTokenClaimsInfo got from the refresh token claims
func (s *DefaultAuthService) getRefreshTokenInfo(refreshToken interface{}) *models.RefreshTokenClaimsInfo {
	claims := s.jwtPrv.GetTokenClaims(refreshToken)

	info := models.RefreshTokenClaimsInfo{
		UserID: s.parseInt32Claim(claims["userId"]),
	}
	return &info
}

func (s *DefaultAuthService) parseStringClaim(value interface{}) string {
	result, _ := value.(string)
	return result
}

func (s *DefaultAuthService) parseInt32Claim(value interface{}) int32 {
	result, _ := value.(float64)
	return int32(result)
}

func (s *DefaultAuthService) parseBoolClaim(value interface{}) bool {
	result, _ := value.(bool)
	return result
}
