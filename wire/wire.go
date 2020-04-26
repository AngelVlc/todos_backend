//+build wireinject

package wire

import (
	"os"

	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitUsersService(db *gorm.DB) services.UsersService {
	wire.Build(CryptoHelperSet, services.NewUsersService)

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
	if len(os.Getenv("TESTING")) > 0 {
		return initMockedAuthService()
	} else {
		return initDefaultAuthService()
	}
}

func initDefaultAuthService() services.AuthService {
	wire.Build(TokenHelperSet, ConfigurationServiceSet, AuthServiceSet)
	return nil
}

func initMockedAuthService() services.AuthService {
	wire.Build(MockedAuthServiceSet)
	return nil
}

func InitConfigurationService() services.ConfigurationService {
	wire.Build(ConfigurationServiceSet)
	return nil
}

var EnvGetterSet = wire.NewSet(
	services.NewOsEnvGetter,
	wire.Bind(new(services.EnvGetter), new(*services.OsEnvGetter)))

var TokenHelperSet = wire.NewSet(
	services.NewJwtTokenHelper,
	wire.Bind(new(services.TokenHelper), new(*services.JwtTokenHelper)))

var MockedTokenHelperSet = wire.NewSet(
	services.NewMockedTokenHelper,
	wire.Bind(new(services.TokenHelper), new(*services.MockedTokenHelper)))

var CryptoHelperSet = wire.NewSet(
	services.NewBcryptHelper,
	wire.Bind(new(services.CryptoHelper), new(*services.BcryptHelper)))

var ConfigurationServiceSet = wire.NewSet(
	EnvGetterSet,
	services.NewDefaultConfigurationService,
	wire.Bind(new(services.ConfigurationService), new(*services.DefaultConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(
	services.NewMockedConfigurationService,
	wire.Bind(new(services.ConfigurationService), new(*services.MockedConfigurationService)))

var AuthServiceSet = wire.NewSet(
	services.NewDefaultAuthService,
	wire.Bind(new(services.AuthService), new(*services.DefaultAuthService)))

var MockedAuthServiceSet = wire.NewSet(
	services.NewMockedAuthService,
	wire.Bind(new(services.AuthService), new(*services.MockedAuthService)))
