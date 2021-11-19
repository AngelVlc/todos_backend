package domain

import "time"

type AuthRepository interface {
	ExistsUser(userName UserName) (bool, error)
	FindUserByName(userName UserName) (*User, error)
	FindUserByID(userID int32) (*User, error)
	GetAllUsers() ([]User, error)
	CreateUser(user *User) error
	DeleteUser(userID int32) error
	UpdateUser(user *User) error

	FindRefreshTokenForUser(refreshToken string, userID int32) (*RefreshToken, error)
	CreateRefreshToken(refreshToken *RefreshToken) error
	DeleteExpiredRefreshTokens(expTime time.Time) error
	GetAllRefreshTokens() ([]RefreshToken, error)
	DeleteRefreshTokensByID(ids []int32) error
}
