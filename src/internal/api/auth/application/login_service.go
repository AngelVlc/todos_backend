package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type LoginService struct {
	repo   domain.AuthRepository
	cfgSvc sharedApp.ConfigurationService
}

func NewLoginService(repo domain.AuthRepository, cfgSvc sharedApp.ConfigurationService) *LoginService {
	return &LoginService{repo, cfgSvc}
}

func (s *LoginService) Login(userName *domain.UserName, password *domain.UserPassword) (*domain.TokenResponse, error) {
	foundUser, err := s.repo.FindUserByName(userName)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	err = foundUser.HasPassword(*password)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid password", InternalError: err}
	}

	tokenSvc := domain.NewTokenService(s.cfgSvc)

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
