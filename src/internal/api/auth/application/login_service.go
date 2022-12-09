package application

import (
	"context"
	"log"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type LoginService struct {
	authRepo  domain.AuthRepository
	usersRepo domain.UsersRepository
	cfgSvr    sharedApp.ConfigurationService
	tokenSrv  domain.TokenService
}

func NewLoginService(authRepo domain.AuthRepository, usersRepo domain.UsersRepository, cfgSvr sharedApp.ConfigurationService, tokenSrv domain.TokenService) *LoginService {
	return &LoginService{authRepo, usersRepo, cfgSvr, tokenSrv}
}

func (s *LoginService) Login(ctx context.Context, userName domain.UserNameValueObject, password domain.UserPassword) (*domain.LoginResponse, error) {
	foundUser, err := s.usersRepo.FindUser(ctx, &domain.User{Name: userName})
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

	refreshTokenExpDate := s.cfgSvr.GetRefreshTokenExpirationTime()

	refreshToken, err := s.tokenSrv.GenerateRefreshToken(foundUser, refreshTokenExpDate)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	txn := newrelic.FromContext(ctx)
	go func(txn *newrelic.Transaction) {
		ctx = newrelic.NewContext(context.Background(), txn)
		defer txn.End()
		err = s.authRepo.CreateRefreshTokenIfNotExist(ctx, &domain.RefreshToken{UserID: foundUser.ID, RefreshToken: refreshToken, ExpirationDate: refreshTokenExpDate})
		if err != nil {
			log.Printf("Error saving the refresh token. Error: %v", err)
		}
	}(txn.NewGoroutine())

	res := domain.LoginResponse{Token: token, RefreshToken: refreshToken, UserID: foundUser.ID, UserName: string(foundUser.Name), IsAdmin: foundUser.IsAdmin}

	return &res, nil
}
