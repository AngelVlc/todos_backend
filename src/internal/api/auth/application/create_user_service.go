package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type CreateUserService struct {
	usersRepo domain.UsersRepository
	passGen   passgen.PasswordGenerator
}

func NewCreateUserService(usersRepo domain.UsersRepository, passGen passgen.PasswordGenerator) *CreateUserService {
	return &CreateUserService{usersRepo, passGen}
}

func (s *CreateUserService) CreateUser(ctx context.Context, userName domain.UserNameValueObject, password domain.UserPassword, isAdmin bool) (*domain.UserRecord, error) {
	err := userName.CheckIfAlreadyExists(ctx, s.usersRepo)
	if err != nil {
		return nil, err
	}

	hasshedPass, err := s.passGen.GenerateFromPassword(string(password))
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user := domain.UserRecord{
		Name:         userName,
		PasswordHash: hasshedPass,
		IsAdmin:      isAdmin,
	}

	err = s.usersRepo.Create(ctx, &user)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user", InternalError: err}
	}

	return &user, nil
}
