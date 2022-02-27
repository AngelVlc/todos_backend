package passgen

import "github.com/stretchr/testify/mock"

type MockedPasswordGenerator struct {
	mock.Mock
}

func NewMockedPasswordGenerator() *MockedPasswordGenerator {
	return &MockedPasswordGenerator{}
}

func (m *MockedPasswordGenerator) GenerateFromPassword(password string) (string, error) {
	args := m.Called(password)

	return args.String(0), args.Error(1)
}
