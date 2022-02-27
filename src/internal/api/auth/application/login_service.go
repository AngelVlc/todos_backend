package application

import (
	"context"
	"log"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type LoginService struct {
	repo     domain.AuthRepository
	cfgSvr   sharedApp.ConfigurationService
	tokenSrv domain.TokenService
}

func NewLoginService(repo domain.AuthRepository, cfgSvr sharedApp.ConfigurationService, tokenSrv domain.TokenService) *LoginService {
	return &LoginService{repo, cfgSvr, tokenSrv}
}

func (s *LoginService) Login(ctx context.Context, userName domain.UserName, password domain.UserPassword) (*domain.LoginResponse, error) {
	foundUser, err := s.repo.FindUserByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	err = foundUser.HasPassword(password)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid password", InternalError: err}
	}

	token, err := s.tokenSrv.GenerateToken(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	refreshTokenExpDate := s.cfgSvr.GetRefreshTokenExpirationDuration()

	refreshToken, err := s.tokenSrv.GenerateRefreshToken(foundUser, refreshTokenExpDate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	txn := newrelic.FromContext(ctx)
	go func(txn *newrelic.Transaction) {
		ctx = newrelic.NewContext(context.Background(), txn)
		defer txn.End()
		err = s.repo.CreateRefreshTokenIfNotExist(ctx, &domain.RefreshToken{UserID: foundUser.ID, RefreshToken: refreshToken, ExpirationDate: refreshTokenExpDate})
		if err != nil {
			log.Printf("Error saving the refresh token. Error: %v", err)
		}
	}(txn.NewGoroutine())

	res := domain.LoginResponse{Token: token, RefreshToken: refreshToken, UserID: foundUser.ID, UserName: string(foundUser.Name), IsAdmin: foundUser.IsAdmin}

	return &res, nil
}
