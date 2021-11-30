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
	CreateUser(user *User) error
	DeleteUser(userID int32) error
	UpdateUser(user *User) error

	FindRefreshTokenForUser(refreshToken string, userID int32) (*RefreshToken, error)
	CreateRefreshTokenIfNotExist(refreshToken *RefreshToken) error
	DeleteExpiredRefreshTokens(expTime time.Time) error
	GetAllRefreshTokens() ([]RefreshToken, error)
	DeleteRefreshTokensByID(ids []int32) error
}
