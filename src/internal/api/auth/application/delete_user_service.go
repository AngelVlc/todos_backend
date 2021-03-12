package application

import (
	"strings"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type DeleteUserService struct {
	repo domain.AuthRepository
}

func NewDeleteUserService(repo domain.AuthRepository) *DeleteUserService {
	return &DeleteUserService{repo}
}

func (s *DeleteUserService) DeleteUser(userID int32) error {
	foundUser, err := s.repo.FindUserByID(userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting user by id", InternalError: err}
	}

	if foundUser == nil {
		return &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	if strings.ToLower(string(foundUser.Name)) == "admin" {
		return &appErrors.BadRequestError{Msg: "It is not possible to delete the admin user"}
	}

	err = s.repo.DeleteUser(userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user", InternalError: err}
	}

	return nil
}
