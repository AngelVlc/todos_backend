package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, query ListRecord) (ListRecord, error)
	ExistsList(ctx context.Context, query ListRecord) (bool, error)
	GetAllLists(ctx context.Context) ([]ListRecord, error)
	GetAllListsForUser(ctx context.Context, userID int32) ([]ListRecord, error)
	CreateList(ctx context.Context, record *ListRecord) error
	DeleteList(ctx context.Context, query ListRecord) error
	UpdateList(ctx context.Context, record *ListRecord) error
	UpdateListItemsCount(ctx context.Context, listID int32) error
}
