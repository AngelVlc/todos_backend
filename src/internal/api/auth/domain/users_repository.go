package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, query UserEntity) (*UserEntity, error)
	ExistsUser(ctx context.Context, query UserEntity) (bool, error)
	GetAll(ctx context.Context) ([]*UserEntity, error)
	Create(ctx context.Context, user *UserEntity) (*UserEntity, error)
	Delete(ctx context.Context, query UserEntity) error
	Update(ctx context.Context, user *UserEntity) (*UserEntity, error)
}
