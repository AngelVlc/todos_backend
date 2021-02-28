package infrastructure

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
