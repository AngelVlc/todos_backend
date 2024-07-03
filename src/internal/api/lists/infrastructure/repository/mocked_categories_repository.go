package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/stretchr/testify/mock"
)

type MockedCategoriesRepository struct {
	mock.Mock
}

func NewMockedCategoriesRepository() *MockedCategoriesRepository {
	return &MockedCategoriesRepository{}
}

func (m *MockedCategoriesRepository) FindCategory(ctx context.Context, query domain.CategoryRecord) (*domain.CategoryRecord, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CategoryRecord), args.Error(1)
}

func (m *MockedCategoriesRepository) ExistsCategory(ctx context.Context, query domain.CategoryRecord) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedCategoriesRepository) GetCategories(ctx context.Context, query domain.CategoryRecord) (domain.CategoryRecords, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(domain.CategoryRecords), args.Error(1)
}

func (m *MockedCategoriesRepository) CreateCategory(ctx context.Context, record *domain.CategoryRecord) error {
	args := m.Called(ctx, record)

	return args.Error(0)
}

func (m *MockedCategoriesRepository) DeleteCategory(ctx context.Context, query domain.CategoryRecord) error {
	args := m.Called(ctx, query)

	return args.Error(0)
}

func (m *MockedCategoriesRepository) UpdateCategory(ctx context.Context, category *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CategoryEntity), args.Error(1)
}
