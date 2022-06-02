package domain

import (
	"context"
	"time"

	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID int32) error
	UpdateUser(ctx context.Context, user *User) error

	FindRefreshTokenForUser(ctx context.Context, refreshToken string, userID int32) (*RefreshToken, error)
	CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *RefreshToken) error
	DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error
	GetAllRefreshTokens(ctx context.Context, paginationInfo *sharedDomain.PaginationInfo) ([]RefreshToken, error)
	DeleteRefreshTokensByID(ctx context.Context, ids []int32) error
}
