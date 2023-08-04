package application

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
)

type GetAllRefreshTokensService struct {
	repo domain.AuthRepository
}

func NewGetAllRefreshTokensService(repo domain.AuthRepository) *GetAllRefreshTokensService {
	return &GetAllRefreshTokensService{repo}
}

func (s *GetAllRefreshTokensService) GetAllRefreshTokens(ctx context.Context, pagInfo *sharedDomain.PaginationInfo) ([]*domain.RefreshTokenEntity, error) {
	found, err := s.repo.GetAllRefreshTokens(ctx, pagInfo)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting refresh tokens", InternalError: err}
	}

	return found, nil
}
