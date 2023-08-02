package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/stretchr/testify/mock"
)

type MockedListsRepository struct {
	mock.Mock
}

func NewMockedListsRepository() *MockedListsRepository {
	return &MockedListsRepository{}
}

func (m *MockedListsRepository) FindList(ctx context.Context, query *domain.ListRecord) (*domain.ListRecord, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.ListRecord), args.Error(1)
}

func (m *MockedListsRepository) ExistsList(ctx context.Context, query *domain.ListRecord) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedListsRepository) GetAllLists(ctx context.Context, userID int32) ([]domain.ListRecord, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.ListRecord), args.Error(1)
}

func (m *MockedListsRepository) CreateList(ctx context.Context, list *domain.ListRecord) error {
	args := m.Called(ctx, list)

	return args.Error(0)
}

func (m *MockedListsRepository) DeleteList(ctx context.Context, list *domain.ListRecord) error {
	args := m.Called(ctx, list)

	return args.Error(0)
}

func (m *MockedListsRepository) UpdateList(ctx context.Context, list *domain.ListRecord) error {
	args := m.Called(ctx, list)

	return args.Error(0)
}

func (m *MockedListsRepository) UpdateListItemsCounter(ctx context.Context, listID int32) error {
	args := m.Called(ctx, listID)

	return args.Error(0)
}
