package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CreateUserService struct {
	authRepo  domain.AuthRepository
	usersRepo domain.UsersRepository
	passGen   passgen.PasswordGenerator
}

func NewCreateUserService(authRepo domain.AuthRepository, usersRepo domain.UsersRepository, passGen passgen.PasswordGenerator) *CreateUserService {
	return &CreateUserService{authRepo, usersRepo, passGen}
}

func (s *CreateUserService) CreateUser(ctx context.Context, userName domain.UserName, password domain.UserPassword, isAdmin bool) (*domain.User, error) {
	err := userName.CheckIfAlreadyExists(ctx, s.usersRepo)
	if err != nil {
		return nil, err
	}

	hasshedPass, err := s.passGen.GenerateFromPassword(string(password))
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user := domain.User{
		Name:         userName,
		PasswordHash: hasshedPass,
		IsAdmin:      isAdmin,
	}

	err = s.authRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user", InternalError: err}
	}

	return &user, nil
}
