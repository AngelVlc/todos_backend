package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedEnvGetter struct {
	mock.Mock
}

func (m *MockedEnvGetter) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestConfigService(t *testing.T) {
	mockedEg := MockedEnvGetter{}
	svc := NewConfigurationService(&mockedEg)

	mockedEg.On("Getenv", "MYSQL_HOST").Return("host")
	mockedEg.On("Getenv", "MYSQL_PORT").Return("port")
	mockedEg.On("Getenv", "MYSQL_USER").Return("user")
	mockedEg.On("Getenv", "MYSQL_PASSWORD").Return("password")
	mockedEg.On("Getenv", "MYSQL_DATABASE").Return("database")

	t.Run("GetDataSource() should return the data source from the env vars when the CLEARDB_DATABASE_URL env var is empty", func(t *testing.T) {
		mockedEg.On("Getenv", "CLEARDB_DATABASE_URL").Return("").Once()
		res := svc.GetDatasource()

		assert.Equal(t, "user:password@(host:port)/database?charset=utf8&parseTime=True&loc=Local", res)

		mockedEg.AssertExpectations(t)
	})

	t.Run("GetDataSource() should return the data source from the CLEARDB_DATABASE_URL env vars when it isn't empty", func(t *testing.T) {
		mockedEg.On("Getenv", "CLEARDB_DATABASE_URL").Return("mysql://user:pass@host/database?reconnect=true").Once()
		res := svc.GetDatasource()

		assert.Equal(t, "user:pass@(host:3306)/database?charset=utf8&parseTime=True&loc=Local", res)

		mockedEg.AssertExpectations(t)
	})

}
