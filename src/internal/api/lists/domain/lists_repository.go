package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, list *ListEntity) (*ListEntity, error)
	ExistsList(ctx context.Context, list *ListEntity) (bool, error)
	GetAllLists(ctx context.Context, userID int32) ([]ListEntity, error)
	CreateList(ctx context.Context, list *ListEntity) error
	DeleteList(ctx context.Context, listID int32, userID int32) error
	UpdateList(ctx context.Context, list *ListEntity) error
	IncrementListCounter(ctx context.Context, listID int32) error
	DecrementListCounter(ctx context.Context, listID int32) error

	FindListItem(ctx context.Context, listItem *ListItemEntity) (*ListItemEntity, error)
	GetAllListItems(ctx context.Context, listID int32, userID int32) ([]ListItemEntity, error)
	CreateListItem(ctx context.Context, listItem *ListItemEntity) error
	DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error
	UpdateListItem(ctx context.Context, listItem *ListItemEntity) error
	BulkUpdateListItems(ctx context.Context, listItems []ListItemEntity) error
	GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error)
}
