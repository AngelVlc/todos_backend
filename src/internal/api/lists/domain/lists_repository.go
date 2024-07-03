package domain

import "context"

type ListsRepository interface {
	/* FindList returns an error if the list doesn't exist */
	FindList(ctx context.Context, query ListRecord) (*ListRecord, error)
	ExistsList(ctx context.Context, query ListRecord) (bool, error)
	GetLists(ctx context.Context, query ListRecord) (ListRecords, error)
	CreateList(ctx context.Context, record *ListRecord) error
	DeleteList(ctx context.Context, query ListRecord) error
	UpdateList(ctx context.Context, record *ListRecord) error
	UpdateListItemsCount(ctx context.Context, listID int32) error
}
