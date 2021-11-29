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
	GetEnvironment() string
	GetHoneyBadgerApiKey() string
	GetNewRelicApiKey() string
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

func (m *MockedConfigurationService) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedConfigurationService) GetHoneyBadgerApiKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedConfigurationService) GetNewRelicApiKey() string {
	args := m.Called()
	return args.String(0)
}

type RealConfigurationService struct{}

func NewRealConfigurationService() *RealConfigurationService {
	return &RealConfigurationService{}
}

func (c *RealConfigurationService) GetDatasource() string {
	host := c.getEnvOrFallback("MYSQL_HOST", "localhost")
	port := c.getEnvOrFallback("MYSQL_PORT", "3306")
	user := c.getEnvOrFallback("MYSQL_USER", "root")
	pass := c.getEnvOrFallback("MYSQL_PASSWORD", "pass")
	dbname := c.getEnvOrFallback("MYSQL_DATABASE", "todos")

	// This env var is set by Heroku
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
	return c.getEnvOrFallback("ADMIN_PASSWORD", "admin")
}

func (c *RealConfigurationService) GetPort() string {
	return c.getEnvOrFallback("PORT", "5001")
}

func (c *RealConfigurationService) GetJwtSecret() string {
	return c.getEnvOrFallback("JWT_SECRET", "mySecret")
}

func (c *RealConfigurationService) GetCorsAllowedOrigins() []string {
	return strings.Split(c.getEnvOrFallback("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ",")
}

func (c *RealConfigurationService) GetTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("TOKEN_EXPIRATION_DURATION", "5m"))
}

func (c *RealConfigurationService) GetRefreshTokenExpirationDate() time.Time {
	return time.Now().Add(c.getDurationEnvVar("REFRESH_TOKEN_EXPIRATION_DURATION", "24h"))
}

func (c *RealConfigurationService) GetEnvironment() string {
	return c.getEnvOrFallback("ENVIRONMENT", "development")
}

func (c *RealConfigurationService) GetHoneyBadgerApiKey() string {
	return c.getEnvOrFallback("HONEYBADGER_API_KEY", "apikey")
}

func (c *RealConfigurationService) GetNewRelicApiKey() string {
	return c.getEnvOrFallback("NEWRELIC_API_KEY", "apikey")
}

func (c *RealConfigurationService) GetDeleteExpiredRefreshTokensInterval() time.Time {
	return time.Now().Add(c.getDurationEnvVar("DELETE_EXPIRED_REFRESH_TOKEN_INTERVAL", "30s"))
}

func (c *RealConfigurationService) getDurationEnvVar(key string, fallback string) time.Duration {
	d, _ := time.ParseDuration(c.getEnvOrFallback(key, fallback))
	return d
}

func (c *RealConfigurationService) getEnvOrFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
