package repositories

import (
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsRepository interface {
	Insert(list *models.List) (int32, error)
	Remove(id int32, userID int32) error
	Update(list *models.List) error
	FindByID(id int32, userID int32) (*models.List, error)
}

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) Insert(list *models.List) (int32, error) {
	args := m.Called(list)
	got := args.Get(0)
	if got == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsRepository) Remove(id int32, userID int32) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockedListsRepository) Update(list *models.List) error {
	args := m.Called(list)
	return args.Error(0)
}

func (m *MockedListsRepository) FindByID(id int32, userID int32) (*models.List, error) {
	args := m.Called(id, userID)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.List), args.Error(1)
}

type DefaultListsRepository struct {
	db *gorm.DB
}

func NewDefaultListsRepository(db *gorm.DB) *DefaultListsRepository {
	return &DefaultListsRepository{db}
}

func (r *DefaultListsRepository) Insert(list *models.List) (int32, error) {
	if err := r.db.Create(&list).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list", InternalError: err}
	}

	return list.ID, nil
}

func (r *DefaultListsRepository) Remove(id int32, userID int32) error {
	if err := r.db.Where(models.List{ID: id, UserID: userID}).Delete(models.List{}).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error deleting user list", InternalError: err}
	}
	return nil
}

func (r *DefaultListsRepository) Update(list *models.List) error {
	if err := r.db.Save(&list).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating list", InternalError: err}
	}

	return nil
}

func (r *DefaultListsRepository) FindByID(id int32, userID int32) (*models.List, error) {
	foundList := models.List{}
	err := r.db.Where(models.List{ID: id, UserID: userID}).Preload("ListItems").Find(&foundList).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user list", InternalError: err}
	}

	return &foundList, nil
}
