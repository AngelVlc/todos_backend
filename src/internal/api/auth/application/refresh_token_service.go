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

func (s *RefreshTokenService) RefreshToken(ctx context.Context, rt string) (string, string, error) {
	parsedRt, err := s.tokenSrv.ParseToken(rt)
	if err != nil {
		return "", "", &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	rtInfo := s.tokenSrv.GetRefreshTokenInfo(parsedRt)

	foundUser, err := s.usersRepo.FindUser(ctx, domain.UserRecord{ID: rtInfo.UserID})
	if err != nil {
		return "", "", err
	}

	entity := foundUser.ToUserEntity()

	token, err := s.tokenSrv.GenerateToken(entity)
	if err != nil {
		return "", "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	refreshTokenExpDate := s.cfgSvr.GetRefreshTokenExpirationTime()

	refreshToken, err := s.tokenSrv.GenerateRefreshToken(entity, refreshTokenExpDate)
	if err != nil {
		return "", "", &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	if err := s.authRepo.ReplaceRefreshToken(ctx, rt, &domain.RefreshTokenEntity{UserID: rtInfo.UserID, RefreshToken: refreshToken, ExpirationDate: refreshTokenExpDate}); err != nil {
		return "", "", &appErrors.UnexpectedError{Msg: "Error replacing refresh token", InternalError: err}
	}

	return token, refreshToken, nil
}
