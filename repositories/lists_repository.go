package repositories

import (
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type ListsRepository interface {
	Insert(userID int32, list *models.List) (int32, error)
}

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) Insert(userID int32, list *models.List) (int32, error) {
	args := m.Called(userID, list)
	got := args.Get(0)
	if got == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

type DefaultListsRepository struct {
	db *gorm.DB
}

func NewDefaultListsRepository(db *gorm.DB) *DefaultListsRepository {
	return &DefaultListsRepository{db}
}

func (r *DefaultListsRepository) Insert(userID int32, list *models.List) (int32, error) {
	if err := r.db.Create(&list).Error; err != nil {
		return 0, &appErrors.UnexpectedError{Msg: "Error inserting list", InternalError: err}
	}

	return list.ID, nil
}
