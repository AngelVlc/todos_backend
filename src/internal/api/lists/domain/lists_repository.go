package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, list *ListRecord) (*ListRecord, error)
	ExistsList(ctx context.Context, list *ListRecord) (bool, error)
	GetAllLists(ctx context.Context, userID int32) ([]ListRecord, error)
	CreateList(ctx context.Context, list *ListRecord) error
	DeleteList(ctx context.Context, listID int32, userID int32) error
	UpdateList(ctx context.Context, list *ListRecord) error
	IncrementListCounter(ctx context.Context, listID int32) error
	DecrementListCounter(ctx context.Context, listID int32) error

	FindListItem(ctx context.Context, listItem *ListItemRecord) (*ListItemRecord, error)
	GetAllListItems(ctx context.Context, listID int32, userID int32) ([]ListItemRecord, error)
	CreateListItem(ctx context.Context, listItem *ListItemRecord) error
	DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error
	UpdateListItem(ctx context.Context, listItem *ListItemRecord) error
	BulkUpdateListItems(ctx context.Context, listItems []ListItemRecord) error
	GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error)
}
