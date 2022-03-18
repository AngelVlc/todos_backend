package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
)

type GetAllUsersService struct {
	repo domain.AuthRepository
}

func NewGetAllUsersService(repo domain.AuthRepository) *GetAllUsersService {
	return &GetAllUsersService{repo}
}

func (s *GetAllUsersService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	foundUsers, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting users", InternalError: err}
	}

	return foundUsers, nil
}
