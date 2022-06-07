package domain

import "context"

type ListsRepository interface {
	FindList(ctx context.Context, list *List) (*List, error)
	GetAllLists(ctx context.Context, userID int32) ([]List, error)
	CreateList(ctx context.Context, list *List) error
	DeleteList(ctx context.Context, listID int32, userID int32) error
	UpdateList(ctx context.Context, list *List) error
	IncrementListCounter(ctx context.Context, listID int32) error
	DecrementListCounter(ctx context.Context, listID int32) error

	FindListItemByID(ctx context.Context, itemID int32, listID int32, userID int32) (*ListItem, error)
	GetAllListItems(ctx context.Context, listID int32, userID int32) ([]ListItem, error)
	CreateListItem(ctx context.Context, listItem *ListItem) error
	DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error
	UpdateListItem(ctx context.Context, listItem *ListItem) error
	BulkUpdateListItems(ctx context.Context, listItems []ListItem) error
	GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error)
}
