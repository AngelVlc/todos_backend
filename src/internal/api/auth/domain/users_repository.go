package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, filter *UserEntity) (*UserEntity, error)
	ExistsUser(ctx context.Context, filter *UserEntity) (bool, error)
	GetAll(ctx context.Context) ([]UserEntity, error)
	Create(ctx context.Context, user *UserEntity) error
	Delete(ctx context.Context, filter *UserEntity) error
	Update(ctx context.Context, user *UserEntity) error
}
