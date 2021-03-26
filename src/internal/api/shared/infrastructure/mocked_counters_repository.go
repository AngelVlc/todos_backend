package infrastructure

import (
	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/stretchr/testify/mock"
)

type MockedCountersRepository struct {
	mock.Mock
}

func NewMockedCountersRepository() *MockedCountersRepository {
	return &MockedCountersRepository{}
}

func (m *MockedCountersRepository) FindByName(name string) (*domain.Counter, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Counter), args.Error(1)
}

func (m *MockedCountersRepository) Create(counter *domain.Counter) error {
	args := m.Called(counter)
	return args.Error(0)
}

func (m *MockedCountersRepository) Update(counter *domain.Counter) error {
	args := m.Called(counter)
	return args.Error(0)
}
