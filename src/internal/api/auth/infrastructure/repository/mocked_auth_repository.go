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

func (m *MockedAuthRepository) FindUserByID(userID int32) (*domain.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedAuthRepository) FindUserByName(userName domain.UserName) (*domain.User, error) {
	args := m.Called(userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedAuthRepository) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockedAuthRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockedAuthRepository) DeleteUser(userID int32) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockedAuthRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockedAuthRepository) FindRefreshTokenForUser(refreshToken string, userID int32) (*domain.RefreshToken, error) {
	args := m.Called(refreshToken, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockedAuthRepository) CreateRefreshToken(refreshToken *domain.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}
