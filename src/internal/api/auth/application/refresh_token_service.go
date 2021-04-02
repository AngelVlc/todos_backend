package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type RefreshTokenService struct {
	repo   domain.AuthRepository
	cfgSvc sharedApp.ConfigurationService
}

func NewRefreshTokenService(repo domain.AuthRepository, cfgSvc sharedApp.ConfigurationService) *LoginService {
	return &LoginService{repo, cfgSvc}
}

func (s *LoginService) RefreshToken(rt string) (*domain.TokenResponse, error) {
	tokenSvc := domain.NewTokenService(s.cfgSvc)

	parsedRt, err := tokenSvc.ParseToken(rt)
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Error parsing the refresh token", InternalError: err}
	}

	if !parsedRt.Valid {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token"}
	}

	rtInfo := tokenSvc.GetRefreshTokenInfo(parsedRt)

	foundUser, err := s.repo.FindUserByID(rtInfo.UserID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user id", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.UnauthorizedError{Msg: "The user no longer exists"}
	}

	token, err := tokenSvc.GenerateToken(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	refreshToken, err := tokenSvc.GenerateRefreshToken(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	res := domain.TokenResponse{Token: token, RefreshToken: refreshToken}

	return &res, nil
}
