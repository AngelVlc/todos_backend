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

	foundUser, err := s.usersRepo.FindUser(ctx, domain.UserRecord{ID: rtInfo.UserID})
	if err != nil {
		return "", err
	}

	if existsRt, err := s.authRepo.ExistsRefreshToken(ctx, domain.RefreshTokenEntity{RefreshToken: rt, UserID: rtInfo.UserID}); err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error getting the refresh token", InternalError: err}
	} else if !existsRt {
		return "", &appErrors.UnauthorizedError{Msg: "The refresh token is not valid"}
	}

	entity := foundUser.ToUserEntity()

	token, err := s.tokenSrv.GenerateToken(entity)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	return token, nil
}
