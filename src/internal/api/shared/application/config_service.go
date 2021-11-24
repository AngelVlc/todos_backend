package application

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stretchr/testify/mock"
)

type ConfigurationService interface {
	GetDatasource() string
	GetAdminPassword() string
	GetPort() string
	GetJwtSecret() string
	GetCorsAllowedOrigins() []string
	GetTokenExpirationDate() time.Time
	GetRefreshTokenExpirationDate() time.Time
}

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

func (m *MockedConfigurationService) GetAdminPassword() string {
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

func (m *MockedConfigurationService) GetTokenExpirationDate() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockedConfigurationService) GetRefreshTokenExpirationDate() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type RealConfigurationService struct{}

func NewRealConfigurationService() *RealConfigurationService {
	return &RealConfigurationService{}
}

func (c *RealConfigurationService) GetDatasource() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")

	clearDbUrl := os.Getenv("CLEARDB_DATABASE_URL")
	if len(clearDbUrl) > 0 {
		clearDbUrl = strings.Replace(clearDbUrl, "mysql://", "", 1)
		clearDbUrl = strings.Replace(clearDbUrl, "?reconnect=true", "", 1)
		parts := strings.Split(clearDbUrl, "@")
		userPass := strings.Split(parts[0], ":")
		user = userPass[0]
		pass = userPass[1]
		hostDbName := strings.Split(parts[1], "/")
		host = hostDbName[0]
		dbname = hostDbName[1]
		port = "3306"
	}

	return fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
}

func (c *RealConfigurationService) GetAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}

func (c *RealConfigurationService) GetPort() string {
	return os.Getenv("PORT")
}

func (c *RealConfigurationService) GetJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}

func (c *RealConfigurationService) GetCorsAllowedOrigins() []string {
	return strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
}

func (c *RealConfigurationService) GetTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("TOKEN_EXPIRATION_DURATION"))
}

func (c *RealConfigurationService) GetRefreshTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("REFRESH_TOKEN_EXPIRATION_DURATION"))
}

func (c *RealConfigurationService) getDurationEnvVar(envVarName string) time.Duration {
	d, _ := time.ParseDuration(os.Getenv(envVarName))
	return d
}
