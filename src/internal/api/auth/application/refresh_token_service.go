package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type RefreshTokenService struct {
	repo     domain.AuthRepository
	cfgSvr   sharedApp.ConfigurationService
	tokenSrv domain.TokenService
}

func NewRefreshTokenService(repo domain.AuthRepository, cfgSvr sharedApp.ConfigurationService, tokenSrv domain.TokenService) *LoginService {
	return &LoginService{repo, cfgSvr, tokenSrv}
}

func (s *LoginService) RefreshToken(rt string) (string, error) {
	parsedRt, err := s.tokenSrv.ParseToken(rt)
	if err != nil {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	rtInfo := s.tokenSrv.GetRefreshTokenInfo(parsedRt)

	foundUser, err := s.repo.FindUserByID(rtInfo.UserID)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error getting user by user id", InternalError: err}
	}

	if foundUser == nil {
		return "", &appErrors.UnauthorizedError{Msg: "The user no longer exists"}
	}

	foundRefreshToken, err := s.repo.FindRefreshTokenForUser(rt, rtInfo.UserID)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error getting the refresh token", InternalError: err}
	}

	if foundRefreshToken == nil {
		return "", &appErrors.UnauthorizedError{Msg: "The refresh token is not valid"}
	}

	token, err := s.tokenSrv.GenerateToken(foundUser)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	return token, nil
}
