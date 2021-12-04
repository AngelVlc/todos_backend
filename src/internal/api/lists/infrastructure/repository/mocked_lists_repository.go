package repository

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/stretchr/testify/mock"
)

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) ExistsList(ctx context.Context, name domain.ListName, userID int32) (bool, error) {
	args := m.Called(ctx, name, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockedListsRepository) FindListByID(ctx context.Context, listID int32, userID int32) (*domain.List, error) {
	args := m.Called(ctx, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.List), args.Error(1)
}

func (m *MockedListsRepository) GetAllLists(ctx context.Context, userID int32) ([]domain.List, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.List), args.Error(1)
}

func (m *MockedListsRepository) CreateList(ctx context.Context, list *domain.List) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockedListsRepository) DeleteList(ctx context.Context, listID int32, userID int32) error {
	args := m.Called(ctx, listID, userID)
	return args.Error(0)
}

func (m *MockedListsRepository) UpdateList(ctx context.Context, list *domain.List) error {
	args := m.Called(ctx, list)
	return args.Error(0)
}

func (m *MockedListsRepository) IncrementListCounter(ctx context.Context, listID int32) error {
	args := m.Called(ctx, listID)
	return args.Error(0)
}

func (m *MockedListsRepository) DecrementListCounter(ctx context.Context, listID int32) error {
	args := m.Called(ctx, listID)
	return args.Error(0)
}

func (m *MockedListsRepository) GetAllListItems(ctx context.Context, listID int32, userID int32) ([]domain.ListItem, error) {
	args := m.Called(ctx, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ListItem), args.Error(1)
}

func (m *MockedListsRepository) FindListItemByID(ctx context.Context, itemID int32, listID int32, userID int32) (*domain.ListItem, error) {
	args := m.Called(ctx, itemID, listID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListItem), args.Error(1)
}

func (m *MockedListsRepository) CreateListItem(ctx context.Context, listItem *domain.ListItem) error {
	args := m.Called(ctx, listItem)
	return args.Error(0)
}

func (m *MockedListsRepository) DeleteListItem(ctx context.Context, itemID int32, listID int32, userID int32) error {
	args := m.Called(ctx, itemID, listID, userID)
	return args.Error(0)
}

func (m *MockedListsRepository) UpdateListItem(ctx context.Context, listItem *domain.ListItem) error {
	args := m.Called(ctx, listItem)
	return args.Error(0)
}

func (m *MockedListsRepository) BulkUpdateListItems(ctx context.Context, listItems []domain.ListItem) error {
	args := m.Called(ctx, listItems)
	return args.Error(0)
}

func (m *MockedListsRepository) GetListItemsMaxPosition(ctx context.Context, listID int32, userID int32) (int32, error) {
	args := m.Called(ctx, listID, userID)
	if args.Get(0) == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}
