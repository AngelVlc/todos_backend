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

func (m *MockedConfigurationService) GetTokenExpirationTime() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func (m *MockedConfigurationService) GetRefreshTokenExpirationTime() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

func (m *MockedConfigurationService) GetDeleteExpiredRefreshTokensIntervalDuration() time.Duration {
	args := m.Called()

	return args.Get(0).(time.Duration)
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

func (m *MockedConfigurationService) GetDomain() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) GetBucketName() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockedConfigurationService) InProduction() bool {
	args := m.Called()

	return args.Bool(0)
}
