package domain

import (
	"context"
	"time"

	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
)

type AuthRepository interface {
	FindRefreshTokenForUser(ctx context.Context, refreshToken string, userID int32) (*RefreshTokenRecord, error)
	CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *RefreshTokenRecord) error
	DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error
	GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]RefreshTokenRecord, error)
	DeleteRefreshTokensByID(ctx context.Context, ids []int32) error
}
