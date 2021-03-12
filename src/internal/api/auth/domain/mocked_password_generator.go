package domain

import "github.com/stretchr/testify/mock"

type MockedPasswordGenerator struct {
	mock.Mock
}

func NewMockedPasswordGenerator() *MockedPasswordGenerator {
	return &MockedPasswordGenerator{}
}

func (m *MockedPasswordGenerator) GenerateFromPassword(password UserPassword) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
