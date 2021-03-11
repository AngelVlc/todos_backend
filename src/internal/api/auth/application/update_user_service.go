package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type UpdateUserService struct {
	repo    domain.AuthRepository
	passGen domain.PasswordGenerator
}

func NewUpdateUserService(repo domain.AuthRepository, passGen domain.PasswordGenerator) *UpdateUserService {
	return &UpdateUserService{repo, passGen}
}

func (s *UpdateUserService) UpdateUser(userID *int32, userName *domain.UserName, password *domain.UserPassword, isAdmin *bool) (*domain.User, error) {
	foundUser, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by id", InternalError: err}
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	if foundUser.IsTheAdminUser() {
		if *userName != domain.UserName("admin") {
			return nil, &appErrors.BadRequestError{Msg: "It is not possible to change the admin user name"}
		}

		if !*isAdmin {
			return nil, &appErrors.BadRequestError{Msg: "The admin user must be an admin"}
		}
	}

	if password != nil {
		hasshedPass, err := s.passGen.GenerateFromPassword(password)
		if err != nil {
			return nil, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
		}

		foundUser.PasswordHash = hasshedPass
	}

	if userName != nil {
		foundUser.Name = *userName
	}

	if isAdmin != nil {
		foundUser.IsAdmin = *isAdmin
	}

	err = s.repo.UpdateUser(foundUser)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user", InternalError: err}
	}

	return s.repo.FindUserByID(userID)
}
