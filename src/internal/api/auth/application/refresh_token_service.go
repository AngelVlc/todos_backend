package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type RefreshTokenService struct {
	authRepo  domain.AuthRepository
	usersRepo domain.UsersRepository
	cfgSvr    sharedApp.ConfigurationService
	tokenSrv  domain.TokenService
}

func NewRefreshTokenService(authRepo domain.AuthRepository, usersRepo domain.UsersRepository, cfgSvr sharedApp.ConfigurationService, tokenSrv domain.TokenService) *RefreshTokenService {
	return &RefreshTokenService{authRepo, usersRepo, cfgSvr, tokenSrv}
}

func (s *RefreshTokenService) RefreshToken(ctx context.Context, rt string) (string, error) {
	parsedRt, err := s.tokenSrv.ParseToken(rt)
	if err != nil {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	rtInfo := s.tokenSrv.GetRefreshTokenInfo(parsedRt)

	foundUser, err := s.usersRepo.FindUser(ctx, domain.UserEntity{ID: rtInfo.UserID})
	if err != nil {
		return "", err
	}

	foundRt, err := s.authRepo.FindRefreshTokenForUser(ctx, rt, rtInfo.UserID)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error getting the refresh token", InternalError: err}
	}

	if foundRt == nil {
		return "", &appErrors.UnauthorizedError{Msg: "The refresh token is not valid"}
	}

	token, err := s.tokenSrv.GenerateToken(foundUser)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	return token, nil
}
