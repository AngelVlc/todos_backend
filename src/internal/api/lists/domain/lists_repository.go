package domain

type ListsRepository interface {
	ExistsList(name ListName, userID int32) (bool, error)
	FindListByID(listID int32, userID int32) (*List, error)
	GetAllLists(userID int32) ([]List, error)
	CreateList(list *List) error
	DeleteList(listID int32, userID int32) error
	UpdateList(list *List) error
	IncrementListCounter(listID int32) error
	DecrementListCounter(listID int32) error

	FindListItemByID(itemID int32, listID int32, userID int32) (*ListItem, error)
	GetAllListItems(listID int32, userID int32) ([]ListItem, error)
	CreateListItem(listItem *ListItem) error
	DeleteListItem(itemID int32, listID int32, userID int32) error
	UpdateListItem(listItem *ListItem) error
	BulkUpdateListItems(listItems []ListItem) error
	GetListItemsMaxPosition(listID int32, userID int32) (int32, error)
}
