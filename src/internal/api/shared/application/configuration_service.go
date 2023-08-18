package application

import "time"

type ConfigurationService interface {
	GetDatasource() string
	GetPort() string
	GetJwtSecret() string
	GetCorsAllowedOrigins() []string
	GetTokenExpirationTime() time.Time
	GetRefreshTokenExpirationTime() time.Time
	GetDeleteExpiredRefreshTokensIntervalDuration() time.Duration
	GetEnvironment() string
	GetHoneyBadgerApiKey() string
	GetNewRelicLicenseKey() string
	GetDomain() string
	GetBucketName() string
	InProduction() bool
	GetAlgoliaAppId() string
	GetAlgoliaApiKey() string
	GetAlgoliaSearchOnlyKey() string
}
