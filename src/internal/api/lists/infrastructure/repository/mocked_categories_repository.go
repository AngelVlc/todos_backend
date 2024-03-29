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

func (m *MockedCategoriesRepository) FindCategory(ctx context.Context, query domain.CategoryEntity) (*domain.CategoryEntity, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CategoryEntity), args.Error(1)
}

func (m *MockedCategoriesRepository) ExistsCategory(ctx context.Context, query domain.CategoryEntity) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedCategoriesRepository) GetAllCategoriesForUser(ctx context.Context, userID int32) ([]*domain.CategoryEntity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.CategoryEntity), args.Error(1)
}

func (m *MockedCategoriesRepository) CreateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	args := m.Called(ctx, list)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CategoryEntity), args.Error(1)
}

func (m *MockedCategoriesRepository) DeleteCategory(ctx context.Context, query domain.CategoryEntity) error {
	args := m.Called(ctx, query)

	return args.Error(0)
}

func (m *MockedCategoriesRepository) UpdateCategory(ctx context.Context, list *domain.CategoryEntity) (*domain.CategoryEntity, error) {
	args := m.Called(ctx, list)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.CategoryEntity), args.Error(1)
}
