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

func (m *MockedUsersRepository) FindUser(ctx context.Context, filter *domain.UserEntity) (*domain.UserEntity, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEntity), args.Error(1)
}

func (m *MockedUsersRepository) ExistsUser(ctx context.Context, filter *domain.UserEntity) (bool, error) {
	args := m.Called(ctx, filter)

	return args.Bool(0), args.Error(1)
}

func (m *MockedUsersRepository) GetAll(ctx context.Context) ([]domain.UserEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.UserEntity), args.Error(1)
}

func (m *MockedUsersRepository) Create(ctx context.Context, user *domain.UserEntity) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

func (m *MockedUsersRepository) Delete(ctx context.Context, filter *domain.UserEntity) error {
	args := m.Called(ctx, filter)

	return args.Error(0)
}

func (m *MockedUsersRepository) Update(ctx context.Context, user *domain.UserEntity) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}
