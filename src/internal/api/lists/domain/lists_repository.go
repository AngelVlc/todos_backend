package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, query ListEntity) (*ListEntity, error)
	ExistsList(ctx context.Context, query ListEntity) (bool, error)
	GetAllLists(ctx context.Context) ([]*ListEntity, error)
	GetAllListsForUser(ctx context.Context, userID int32) ([]*ListEntity, error)
	CreateList(ctx context.Context, list *ListEntity) (*ListEntity, error)
	DeleteList(ctx context.Context, query ListEntity) error
	UpdateList(ctx context.Context, list *ListEntity) (*ListEntity, error)
	UpdateListItemsCount(ctx context.Context, listID int32) error
}
