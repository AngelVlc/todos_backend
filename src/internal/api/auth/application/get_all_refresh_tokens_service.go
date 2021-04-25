package application

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type GetAllRefreshTokensService struct {
	repo domain.AuthRepository
}

func NewGetAllRefreshTokensService(repo domain.AuthRepository) *GetAllUsersService {
	return &GetAllUsersService{repo}
}

func (s *GetAllUsersService) GetAllRefreshTokens() ([]domain.RefreshToken, error) {
	found, err := s.repo.GetAllRefreshTokens()
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting refresh tokens", InternalError: err}
	}

	return found, nil
}
