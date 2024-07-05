package domain

import "context"

type UsersRepository interface {
	FindUser(ctx context.Context, query UserRecord) (*UserRecord, error)
	ExistsUser(ctx context.Context, query UserRecord) (bool, error)
	GetAll(ctx context.Context) ([]*UserEntity, error)
	Create(ctx context.Context, record *UserRecord) error
	Delete(ctx context.Context, query UserRecord) error
	Update(ctx context.Context, record *UserRecord) error
}
