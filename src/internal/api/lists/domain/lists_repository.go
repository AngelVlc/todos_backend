package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, query *ListRecord) (*ListRecord, error)
	ExistsList(ctx context.Context, query *ListRecord) (bool, error)
	GetAllLists(ctx context.Context, userID int32) ([]ListRecord, error)
	CreateList(ctx context.Context, list *ListRecord) error
	DeleteList(ctx context.Context, list *ListRecord) error
	UpdateList(ctx context.Context, list *ListRecord) error
	UpdateListItemsCounter(ctx context.Context, listID int32) error
}
