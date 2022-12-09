package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetAllUsersService struct {
	repo domain.UsersRepository
}

func NewGetAllUsersService(repo domain.UsersRepository) *GetAllUsersService {
	return &GetAllUsersService{repo}
}

func (s *GetAllUsersService) GetAllUsers(ctx context.Context) ([]domain.UserEntity, error) {
	foundUsers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting users", InternalError: err}
	}

	return foundUsers, nil
}
