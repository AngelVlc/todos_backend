package application

import (
	"context"
	"strings"

	"github.com/AngelVlc/todos_backend/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
)

type DeleteUserService struct {
	repo domain.AuthRepository
}

func NewDeleteUserService(repo domain.AuthRepository) *DeleteUserService {
	return &DeleteUserService{repo}
}

func (s *DeleteUserService) DeleteUser(ctx context.Context, userID int32) error {
	foundUser, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if strings.ToLower(string(foundUser.Name)) == "admin" {
		return &appErrors.BadRequestError{Msg: "It is not possible to delete the admin user"}
	}

	err = s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the user", InternalError: err}
	}

	return nil
}
