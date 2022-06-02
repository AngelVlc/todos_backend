package application

import (
	"context"
	"strings"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type DeleteUserService struct {
	authRepo  domain.AuthRepository
	usersRepo domain.UsersRepository
}

func NewDeleteUserService(authRepo domain.AuthRepository, usersRepo domain.UsersRepository) *DeleteUserService {
	return &DeleteUserService{authRepo, usersRepo}
}

func (s *DeleteUserService) DeleteUser(ctx context.Context, userID int32) error {
	foundUser, err := s.usersRepo.FindUser(ctx, &domain.User{ID: userID})
	if err != nil {
		return err
	}

	if strings.ToLower(string(foundUser.Name)) == "admin" {
		return &appErrors.BadRequestError{Msg: "It is not possible to delete the admin user"}
	}

	err = s.authRepo.DeleteUser(ctx, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user", InternalError: err}
	}

	return nil
}
