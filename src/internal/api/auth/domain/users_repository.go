package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, filter *UserRecord) (*UserRecord, error)
	ExistsUser(ctx context.Context, filter *UserRecord) (bool, error)
	GetAll(ctx context.Context) ([]UserRecord, error)
	Create(ctx context.Context, user *UserRecord) error
	Delete(ctx context.Context, filter *UserRecord) error
	Update(ctx context.Context, user *UserRecord) error
}
