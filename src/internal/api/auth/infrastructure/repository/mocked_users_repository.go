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

func (m *MockedUsersRepository) FindUser(ctx context.Context, query *domain.UserRecord) (*domain.UserRecord, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserRecord), args.Error(1)
}

func (m *MockedUsersRepository) ExistsUser(ctx context.Context, query *domain.UserRecord) (bool, error) {
	args := m.Called(ctx, query)

	return args.Bool(0), args.Error(1)
}

func (m *MockedUsersRepository) GetAll(ctx context.Context) ([]domain.UserRecord, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.UserRecord), args.Error(1)
}

func (m *MockedUsersRepository) Create(ctx context.Context, user *domain.UserRecord) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

func (m *MockedUsersRepository) Delete(ctx context.Context, query *domain.UserRecord) error {
	args := m.Called(ctx, query)

	return args.Error(0)
}

func (m *MockedUsersRepository) Update(ctx context.Context, user *domain.UserRecord) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}
