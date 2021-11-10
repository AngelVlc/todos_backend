package repository

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/stretchr/testify/mock"
)

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) FindListByID(listID int32, userID int32) (*domain.List, error) {
	args := m.Called(listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.List), args.Error(1)
}

func (m *MockedListsRepository) FindListByName(name domain.ListName, userID int32) (*domain.List, error) {
	args := m.Called(name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.List), args.Error(1)
}

func (m *MockedListsRepository) GetAllLists(userID int32) ([]domain.List, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.List), args.Error(1)
}

func (m *MockedListsRepository) CreateList(list *domain.List) error {
	args := m.Called(list)
	return args.Error(0)
}

func (m *MockedListsRepository) DeleteList(listID int32, userID int32) error {
	args := m.Called(listID, userID)
	return args.Error(0)
}

func (m *MockedListsRepository) UpdateList(list *domain.List) error {
	args := m.Called(list)
	return args.Error(0)
}

func (m *MockedListsRepository) IncrementListCounter(listID int32) error {
	args := m.Called(listID)
	return args.Error(0)
}

func (m *MockedListsRepository) DecrementListCounter(listID int32) error {
	args := m.Called(listID)
	return args.Error(0)
}

func (m *MockedListsRepository) GetAllListItems(listID int32, userID int32) ([]domain.ListItem, error) {
	args := m.Called(listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ListItem), args.Error(1)
}

func (m *MockedListsRepository) FindListItemByID(itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	args := m.Called(itemID, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListItem), args.Error(1)
}

func (m *MockedListsRepository) CreateListItem(listItem *domain.ListItem) error {
	args := m.Called(listItem)
	return args.Error(0)
}

func (m *MockedListsRepository) DeleteListItem(itemID int32, listID int32, userID int32) error {
	args := m.Called(itemID, listID, userID)
	return args.Error(0)
}

func (m *MockedListsRepository) UpdateListItem(listItem *domain.ListItem) error {
	args := m.Called(listItem)
	return args.Error(0)
}

func (m *MockedListsRepository) BulkUpdateListItems(listItems []domain.ListItem) error {
	args := m.Called(listItems)
	return args.Error(0)
}

func (m *MockedListsRepository) GetListItemsMaxPosition(listID int32, userID int32) (int32, error) {
	args := m.Called(listID, userID)
	if args.Get(0) == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}
