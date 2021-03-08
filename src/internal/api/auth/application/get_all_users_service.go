package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
)

type GetAllUsersService struct {
	repo domain.AuthRepository
}

func NewGetAllUsersService(repo domain.AuthRepository) *GetAllUsersService {
	return &GetAllUsersService{repo}
}

func (s *GetAllUsersService) GetAllUsers() ([]*domain.AuthUser, error) {
	foundUsers, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting users", InternalError: err}
	}

	return foundUsers, nil
}
