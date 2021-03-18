package domain

type ListsRepository interface {
	FindListByID(listID int32, userID int32) (*List, error)
	GetAllLists(userID int32) ([]List, error)
	CreateList(list *List) error
	DeleteList(listID int32, userID int32) error
	UpdateList(list *List) error

	FindListItemByID(itemID int32, listID int32, userID int32) (*ListItem, error)
	GetAllListItems(listID int32, userID int32) ([]ListItem, error)
	CreateListItem(listItem *ListItem) error
	DeleteListItem(itemID int32, listID int32, userID int32) error
	UpdateListItem(listItem *ListItem) error
}
