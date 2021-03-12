package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type CreateUserService struct {
	repo    domain.AuthRepository
	passGen domain.PasswordGenerator
}

func NewCreateUserService(repo domain.AuthRepository, passGen domain.PasswordGenerator) *CreateUserService {
	return &CreateUserService{repo, passGen}
}

func (s *CreateUserService) CreateUser(userName domain.UserName, password domain.UserPassword, isAdmin bool) (*domain.User, error) {
	foundUser, err := s.repo.FindUserByName(userName)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	if foundUser != nil {
		return nil, &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	hasshedPass, err := s.passGen.GenerateFromPassword(password)
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
