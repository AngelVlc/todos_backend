package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UpdateUserService struct {
	repo    domain.AuthRepository
	passGen passgen.PasswordGenerator
}

func NewUpdateUserService(repo domain.AuthRepository, passGen passgen.PasswordGenerator) *UpdateUserService {
	return &UpdateUserService{repo, passGen}
}

func (s *UpdateUserService) UpdateUser(userID int32, userName domain.UserName, password domain.UserPassword, isAdmin bool) (*domain.User, error) {
	foundUser, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by id", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	if foundUser.IsTheAdminUser() {
		if userName != domain.UserName("admin") {
			return nil, &appErrors.BadRequestError{Msg: "It is not possible to change the admin user name"}
		}

		if !isAdmin {
			return nil, &appErrors.BadRequestError{Msg: "The admin user must be an admin"}
		}
	}

	if foundUser.Name != userName {
		err = userName.CheckIfAlreadyExists(s.repo)
		if err != nil {
			return nil, err
		}
	}

	if len(password) > 0 {
		hasshedPass, err := s.passGen.GenerateFromPassword(string(password))
		if err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
		}

		foundUser.PasswordHash = hasshedPass
	}

	foundUser.Name = userName
	foundUser.IsAdmin = isAdmin

	err = s.repo.UpdateUser(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user", InternalError: err}
	}

	return foundUser, nil
}
