package repositories

import (
	"github.com/AngelVlc/todos/internal/api/models"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsRepository interface {
	Create(list *models.List) (int32, error)
	Delete(id int32, userID int32) error
	Update(list *models.List) error
	FindByID(id int32, userID int32) (*models.List, error)
	GetAll(userID int32) ([]*models.List, error)
}

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) Create(list *models.List) (int32, error) {
	args := m.Called(list)
	got := args.Get(0)
	if got == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedListsRepository) Delete(id int32, userID int32) error {
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

func (m *MockedListsRepository) GetAll(userID int32) ([]*models.List, error) {
	args := m.Called(userID)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.List), args.Error(1)
}

type DefaultListsRepository struct {
	db *gorm.DB
}

func NewDefaultListsRepository(db *gorm.DB) *DefaultListsRepository {
	return &DefaultListsRepository{db}
}

func (r *DefaultListsRepository) Create(list *models.List) (int32, error) {
	if err := r.db.Create(list).Error; err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error creating list", InternalError: err}
	}

	return list.ID, nil
}

func (r *DefaultListsRepository) Delete(id int32, userID int32) error {
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

func (r *DefaultListsRepository) GetAll(userID int32) ([]*models.List, error) {
	res := []*models.List{}
	if err := r.db.Where(models.List{UserID: userID}).Select("id,name").Find(&res).Error; err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user lists", InternalError: err}
	}
	return res, nil
}
