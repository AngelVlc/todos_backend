package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, filter *User) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, filter *User) error
	Update(ctx context.Context, user *User) error
}
