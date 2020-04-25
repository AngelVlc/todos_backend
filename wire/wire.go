//+build wireinject

package wire

import (
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitUsersService(db *gorm.DB) services.UsersService {
	wire.Build(CryptoProviderSet, services.NewUsersService)

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
	wire.Build(TokenProviderSet, EnvGetterSet, services.NewConfigurationService, AuthServiceSet)
	return &services.DefaultAuthService{}
}

func InitConfigurationService() services.ConfigurationService {
	wire.Build(EnvGetterSet, services.NewConfigurationService)
	return services.ConfigurationService{}
}

var EnvGetterSet = wire.NewSet(
	services.NewOsEnvGetter,
	wire.Bind(new(services.EnvGetter), new(*services.OsEnvGetter)))

var TokenProviderSet = wire.NewSet(
	services.NewJwtTokenHelper,
	wire.Bind(new(services.TokenHelper), new(*services.JwtTokenHelper)))

var CryptoProviderSet = wire.NewSet(
	services.NewBcryptHelper,
	wire.Bind(new(services.CryptoHelper), new(*services.BcryptHelper)))

var AuthServiceSet = wire.NewSet(
	services.NewDefaultAuthService,
	wire.Bind(new(services.AuthService), new(*services.DefaultAuthService)))
