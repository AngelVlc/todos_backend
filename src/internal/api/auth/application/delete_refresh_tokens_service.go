package application

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type DeleteRefreshTokensService struct {
	repo domain.AuthRepository
}

func NewDeleteRefreshTokensService(repo domain.AuthRepository) *DeleteRefreshTokensService {
	return &DeleteRefreshTokensService{repo}
}

func (s *DeleteRefreshTokensService) DeleteRefreshTokens(ctx context.Context, ids []int32) error {
	err := s.repo.DeleteRefreshTokensByID(ctx, ids)
	if err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting the refresh tokens", InternalError: err}
	}

	return nil
}
