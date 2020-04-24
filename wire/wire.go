//+build wireinject

package wire

import (
	"github.com/AngelVlc/todos/providers"
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitUsersService(db *gorm.DB) services.UsersService {
	wire.Build(services.NewUsersService)

	return services.UsersService{}
}

func InitCountersService(db *gorm.DB) services.CountersService {
	wire.Build(services.NewCountersService)

	return services.CountersService{}
}

func InitListsService(db *gorm.DB) services.ListsService {
	wire.Build(services.NewListsService)

	return services.ListsService{}
}

func InitAuthService() services.AuthService {
	wire.Build(TokenProviderSet, ConfigurationServiceSet, services.NewAuthService)

	return services.AuthService{}
}

func InitConfigurationService() services.ConfigurationService {
	wire.Build(ConfigurationServiceSet)
	return services.ConfigurationService{}
}

var ConfigurationServiceSet = wire.NewSet(
	providers.NewOsEnvGetter,
	wire.Bind(new(providers.EnvGetter), new(*providers.OsEnvGetter)),
	services.NewConfigurationService)

var TokenProviderSet = wire.NewSet(
	providers.NewJwtTokenProvider,
	wire.Bind(new(providers.TokenProvider), new(*providers.JwtTokenProvider)))
