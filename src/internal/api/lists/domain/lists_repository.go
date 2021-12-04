package domain

import "context"

type ListsRepository interface {
	ExistsList(ctx context.Context, name ListName, userID int32) (bool, error)
	FindListByID(ctx context.Context, listID int32, userID int32) (*List, error)
	GetAllLists(ctx context.Context, userID int32) ([]List, error)
	CreateList(ctx context.Context, list *List) error
	DeleteList(ctx context.Context, listID int32, userID int32) error
	UpdateList(ctx context.Context, list *List) error
	IncrementListCounter(ctx context.Context, listID int32) error
	DecrementListCounter(ctx context.Context, listID int32) error

	FindListItemByID(itemID int32, listID int32, userID int32) (*ListItem, error)
	GetAllListItems(listID int32, userID int32) ([]ListItem, error)
	CreateListItem(listItem *ListItem) error
	DeleteListItem(itemID int32, listID int32, userID int32) error
	UpdateListItem(listItem *ListItem) error
	BulkUpdateListItems(listItems []ListItem) error
	GetListItemsMaxPosition(listID int32, userID int32) (int32, error)
}
