package application

import (
	"context"
	"strings"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type DeleteUserService struct {
	usersRepo domain.UsersRepository
}

func NewDeleteUserService(usersRepo domain.UsersRepository) *DeleteUserService {
	return &DeleteUserService{usersRepo}
}

func (s *DeleteUserService) DeleteUser(ctx context.Context, userID int32) error {
	foundUser, err := s.usersRepo.FindUser(ctx, domain.UserEntity{ID: userID})
	if err != nil {
		return err
	}

	if strings.ToLower(string(foundUser.Name.String())) == "admin" {
		return &appErrors.BadRequestError{Msg: "It is not possible to delete the admin user"}
	}

	err = s.usersRepo.Delete(ctx, domain.UserEntity{ID: userID})
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user", InternalError: err}
	}

	return nil
}
