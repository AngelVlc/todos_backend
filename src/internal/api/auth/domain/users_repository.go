package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, user *User) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user *User) error
}
