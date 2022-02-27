package application

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockedConfigurationService struct {
	mock.Mock
}

func NewMockedConfigurationService() *MockedConfigurationService {
	return &MockedConfigurationService{}
}

func (m *MockedConfigurationService) GetDatasource() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) GetPort() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) GetJwtSecret() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) GetCorsAllowedOrigins() []string {
	args := m.Called()

	return args.Get(0).([]string)
}

func (m *MockedConfigurationService) GetTokenExpirationDuration() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func (m *MockedConfigurationService) GetRefreshTokenExpirationDuration() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func (m *MockedConfigurationService) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedConfigurationService) GetHoneyBadgerApiKey() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) GetNewRelicLicenseKey() string {
	args := m.Called()

	return args.String(0)
}
