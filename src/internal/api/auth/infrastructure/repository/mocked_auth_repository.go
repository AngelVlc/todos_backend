package repository

import (
	"context"
	"sync"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/stretchr/testify/mock"
)

type MockedAuthRepository struct {
	mock.Mock
	mu sync.Mutex
	Wg sync.WaitGroup
}

func NewMockedAuthRepository() *MockedAuthRepository {
	return &MockedAuthRepository{}
}

func (m *MockedAuthRepository) ExistsUser(ctx context.Context, userName domain.UserName) (bool, error) {
	args := m.Called(ctx, userName)
	return args.Bool(0), args.Error(1)
}

func (m *MockedAuthRepository) FindUserByID(ctx context.Context, userID int32) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedAuthRepository) FindUserByName(ctx context.Context, userName domain.UserName) (*domain.User, error) {
	args := m.Called(ctx, userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedAuthRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockedAuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockedAuthRepository) DeleteUser(ctx context.Context, userID int32) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockedAuthRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockedAuthRepository) FindRefreshTokenForUser(refreshToken string, userID int32) (*domain.RefreshToken, error) {
	args := m.Called(refreshToken, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockedAuthRepository) CreateRefreshTokenIfNotExist(refreshToken *domain.RefreshToken) error {
	defer m.Wg.Done()

	m.mu.Lock()
	defer m.mu.Unlock()

	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockedAuthRepository) DeleteExpiredRefreshTokens(expTime time.Time) error {
	args := m.Called(expTime)
	return args.Error(0)
}

func (m *MockedAuthRepository) GetAllRefreshTokens() ([]domain.RefreshToken, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefreshToken), args.Error(1)
}

func (m *MockedAuthRepository) DeleteRefreshTokensByID(ids []int32) error {
	args := m.Called(ids)
	return args.Error(0)
}
