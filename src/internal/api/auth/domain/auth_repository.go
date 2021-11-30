package domain

import (
	"context"
	"time"
)

type AuthRepository interface {
	ExistsUser(ctx context.Context, userName UserName) (bool, error)
	FindUserByName(ctx context.Context, userName UserName) (*User, error)
	FindUserByID(ctx context.Context, userID int32) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID int32) error
	UpdateUser(ctx context.Context, user *User) error

	FindRefreshTokenForUser(ctx context.Context, refreshToken string, userID int32) (*RefreshToken, error)
	CreateRefreshTokenIfNotExist(ctx context.Context, refreshToken *RefreshToken) error
	DeleteExpiredRefreshTokens(ctx context.Context, expTime time.Time) error
	GetAllRefreshTokens() ([]RefreshToken, error)
	DeleteRefreshTokensByID(ids []int32) error
}
