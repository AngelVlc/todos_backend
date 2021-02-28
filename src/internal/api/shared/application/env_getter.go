package application

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type EnvGetter interface {
	Getenv(key string) string
}

type MockedEnvGetter struct {
	mock.Mock
}

func (m *MockedEnvGetter) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

type OsEnvGetter struct{}

func NewOsEnvGetter() *OsEnvGetter {
	return new(OsEnvGetter)
}

func (b *OsEnvGetter) Getenv(key string) string {
	return os.Getenv(key)
}
