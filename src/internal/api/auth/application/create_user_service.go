package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type CreateUserService struct {
	repo    domain.AuthRepository
	passGen passgen.PasswordGenerator
}

func NewCreateUserService(repo domain.AuthRepository, passGen passgen.PasswordGenerator) *CreateUserService {
	return &CreateUserService{repo, passGen}
}

func (s *CreateUserService) CreateUser(userName domain.UserName, password domain.UserPassword, isAdmin bool) (*domain.User, error) {
	err := userName.CheckIfAlreadyExists(s.repo)
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

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating the user", InternalError: err}
	}

	return &user, nil
}
