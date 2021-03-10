package repository

import (
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/stretchr/testify/mock"
)

type MockedAuthRepository struct {
	mock.Mock
}

func NewMockedAuthRepository() *MockedAuthRepository {
	return &MockedAuthRepository{}
}

func (m *MockedAuthRepository) FindUserByID(userID *int32) (*domain.AuthUser, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthUser), args.Error(1)
}

func (m *MockedAuthRepository) FindUserByName(userName *domain.AuthUserName) (*domain.AuthUser, error) {
	args := m.Called(userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthUser), args.Error(1)
}

func (m *MockedAuthRepository) GetAllUsers() ([]*domain.AuthUser, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AuthUser), args.Error(1)
}

func (m *MockedAuthRepository) CreateUser(user *domain.AuthUser) (int32, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedAuthRepository) DeleteUser(userID *int32) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockedAuthRepository) UpdateUser(user *domain.AuthUser) error {
	args := m.Called(user)
	return args.Error(0)
}
