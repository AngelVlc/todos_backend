package repository

import (
	"context"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/stretchr/testify/mock"
)

type MockedUsersRepository struct {
	mock.Mock
}

func NewMockedUsersRepository() *MockedUsersRepository {
	return &MockedUsersRepository{}
}

func (m *MockedUsersRepository) FindUser(ctx context.Context, query domain.UserEntity) (*domain.UserEntity, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEntity), args.Error(1)
}

func (m *MockedUsersRepository) ExistsUser(ctx context.Context, query domain.UserEntity) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedUsersRepository) GetAll(ctx context.Context) ([]*domain.UserEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.UserEntity), args.Error(1)
}

func (m *MockedUsersRepository) Create(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.UserEntity), args.Error(1)
}

func (m *MockedUsersRepository) Delete(ctx context.Context, query domain.UserEntity) error {
	args := m.Called(ctx, query)

	return args.Error(0)
}

func (m *MockedUsersRepository) Update(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.UserEntity), args.Error(1)
}
