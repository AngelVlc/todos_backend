package application

import (
	"fmt"
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

type RealConfigurationService struct {
	eg EnvGetter
}

func NewRealConfigurationService(eg EnvGetter) *RealConfigurationService {
	return &RealConfigurationService{eg}
}

func (c *RealConfigurationService) GetDatasource() string {
	host := c.eg.Getenv("MYSQL_HOST")
	port := c.eg.Getenv("MYSQL_PORT")
	user := c.eg.Getenv("MYSQL_USER")
	pass := c.eg.Getenv("MYSQL_PASSWORD")
	dbname := c.eg.Getenv("MYSQL_DATABASE")

	clearDbUrl := c.eg.Getenv("CLEARDB_DATABASE_URL")
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
	return c.eg.Getenv("ADMIN_PASSWORD")
}

func (c *RealConfigurationService) GetPort() string {
	return c.eg.Getenv("PORT")
}

func (c *RealConfigurationService) GetJwtSecret() string {
	return c.eg.Getenv("JWT_SECRET")
}

func (c *RealConfigurationService) GetCorsAllowedOrigins() []string {
	return strings.Split(c.eg.Getenv("CORS_ALLOWED_ORIGINS"), ",")
}

func (c *RealConfigurationService) GetTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("TOKEN_EXPIRATION_IN_SECONDS"))
}

func (c *RealConfigurationService) GetRefreshTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("REFRESH_TOKEN_EXPIRATION_IN_SECONDS"))
}

func (c *RealConfigurationService) getDurationEnvVar(envVarName string) time.Duration {
	d, _ := time.ParseDuration(c.eg.Getenv(envVarName))
	return d
}
