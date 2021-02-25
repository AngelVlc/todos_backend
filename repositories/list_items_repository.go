package repositories

import (
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListItemsRepository interface {
	FindByID(id int32, listID int32, userID int32) (*models.ListItem, error)
	Insert(item *models.ListItem) (int32, error)
}

type MockedListItemsRepository struct {
	mock.Mock
}

func NewMockedListItemsRepository() *MockedListItemsRepository {
	return &MockedListItemsRepository{}
}

func (m *MockedListItemsRepository) FindByID(id int32, listID int32, userID int32) (*models.ListItem, error) {
	args := m.Called(id, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ListItem), args.Error(1)
}

func (m *MockedListItemsRepository) Insert(item *models.ListItem) (int32, error) {
	args := m.Called(item)
	got := args.Get(0)
	if got == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

type DefaultListItemsRepository struct {
	db *gorm.DB
}

func NewDefaultListItemsRepository(db *gorm.DB) *DefaultListItemsRepository {
	return &DefaultListItemsRepository{db}
}

func (r *DefaultListItemsRepository) FindByID(id int32, listID int32, userID int32) (*models.ListItem, error) {
	found := models.ListItem{}
	err := r.db.Joins("JOIN lists on listItems.listId=lists.id").Where(models.List{ID: listID, UserID: userID}).Where(models.ListItem{ID: id}).Find(&found).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user list item", InternalError: err}
	}

	return &found, nil
}

func (r *DefaultListItemsRepository) Insert(item *models.ListItem) (int32, error) {
	if err := r.db.Create(item).Error; err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error inserting list item", InternalError: err}
	}

	return item.ID, nil
}
