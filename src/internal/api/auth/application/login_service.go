package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type LoginService struct {
	repo     domain.AuthRepository
	cfgSvr   sharedApp.ConfigurationService
	tokenSrv domain.TokenService
}

func NewLoginService(repo domain.AuthRepository, cfgSvr sharedApp.ConfigurationService, tokenSrv domain.TokenService) *LoginService {
	return &LoginService{repo, cfgSvr, tokenSrv}
}

func (s *LoginService) Login(userName domain.UserName, password domain.UserPassword) (*domain.TokenResponse, error) {
	foundUser, err := s.repo.FindUserByName(userName)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	err = foundUser.HasPassword(password)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid password", InternalError: err}
	}

	token, err := s.tokenSrv.GenerateToken(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	refreshTokenExpDate := s.cfgSvr.GetRefreshTokenExpirationDate()

	refreshToken, err := s.tokenSrv.GenerateRefreshToken(foundUser, refreshTokenExpDate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	err = s.repo.CreateRefreshToken(&domain.RefreshToken{UserID: foundUser.ID, RefreshToken: refreshToken, ExpirationDate: refreshTokenExpDate})
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error saving the refresh token", InternalError: err}
	}

	res := domain.TokenResponse{Token: token, RefreshToken: refreshToken}

	return &res, nil
}
