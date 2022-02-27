package application

import "time"

type ConfigurationService interface {
	GetDatasource() string
	GetPort() string
	GetJwtSecret() string
	GetCorsAllowedOrigins() []string
	GetTokenExpirationDuration() time.Time
	GetRefreshTokenExpirationDuration() time.Time
	GetDeleteExpiredTokensIntervalTime() time.Duration
	GetEnvironment() string
	GetHoneyBadgerApiKey() string
	GetNewRelicLicenseKey() string
}
