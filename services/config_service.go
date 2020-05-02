package services

import (
	"fmt"
	"strings"

	"github.com/stretchr/testify/mock"
)

type ConfigurationService interface {
	GetDatasource() string
	GetAdminPassword() string
	GetPort() string
	GetJwtSecret() string
	GetCorsAllowedOrigins() []string
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

type DefaultConfigurationService struct {
	eg EnvGetter
}

func NewDefaultConfigurationService(eg EnvGetter) *DefaultConfigurationService {
	return &DefaultConfigurationService{eg}
}

func (c *DefaultConfigurationService) GetDatasource() string {
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

func (c *DefaultConfigurationService) GetAdminPassword() string {
	return c.eg.Getenv("ADMIN_PASSWORD")
}

func (c *DefaultConfigurationService) GetPort() string {
	return c.eg.Getenv("PORT")
}

func (c *DefaultConfigurationService) GetJwtSecret() string {
	return c.eg.Getenv("JWT_SECRET")
}

func (c *DefaultConfigurationService) GetCorsAllowedOrigins() []string {
	return strings.Split(c.eg.Getenv("CORS_ALLOWED_ORIGINS"), ",")
}
