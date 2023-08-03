package domain

import (
	"context"
	"time"

	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
)

type AuthRepository interface {
	ExistsRefreshToken(ctx context.Context, query RefreshTokenEntity) (bool, error)
	CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *RefreshTokenEntity) error
	DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error
	GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]*RefreshTokenEntity, error)
	DeleteRefreshTokensByID(ctx context.Context, ids []int32) error
}
