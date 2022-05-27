package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, user *User) (*User, error)
}
