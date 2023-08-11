package application

import (
	"fmt"
	"os"
	"strings"
	"time"
)

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
	tls := c.getEnvOrFallback("MYSQL_TLS", "false")
	options := fmt.Sprintf("charset=utf8&parseTime=True&loc=Local&tls=%v", tls)

	return fmt.Sprintf("%v:%v@(%v:%v)/%v?%v", user, pass, host, port, dbname, options)
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

func (c *RealConfigurationService) GetTokenExpirationTime() time.Time {
	return time.Now().Add(c.getDurationEnvVar("TOKEN_EXPIRATION_TIME", "5m"))
}

func (c *RealConfigurationService) GetRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(c.getDurationEnvVar("REFRESH_TOKEN_EXPIRATION_TIME", "24h"))
}

func (c *RealConfigurationService) GetDeleteExpiredRefreshTokensIntervalDuration() time.Duration {
	return c.getDurationEnvVar("DELETE_EXPIRED_REFRESH_TOKEN_INTERVAL", "30s")
}

func (c *RealConfigurationService) GetEnvironment() string {
	return c.getEnvOrFallback("ENVIRONMENT", "development")
}

func (c *RealConfigurationService) GetHoneyBadgerApiKey() string {
	return c.getEnvOrFallback("HONEYBADGER_API_KEY", "apikey")
}

func (c *RealConfigurationService) GetNewRelicLicenseKey() string {
	return c.getEnvOrFallback("NEW_RELIC_LICENSE_KEY", "apikey")
}

func (c *RealConfigurationService) GetDomain() string {
	return c.getEnvOrFallback("DOMAIN", "domain")
}

func (c *RealConfigurationService) GetBucketName() string {
	return c.getEnvOrFallback("BUCKET_NAME", "todos-backend")
}

func (c *RealConfigurationService) InProduction() bool {
	return c.GetEnvironment() == "production"
}

func (c *RealConfigurationService) GetAlgoliaAppId() string {
	return c.getEnvOrFallback("ALGOLIA_APP_ID", "algolia-app-id")
}

func (c *RealConfigurationService) GetAlgoliaApiKey() string {
	return c.getEnvOrFallback("ALGOLIA_API_KEY", "algolia-api-key")
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
