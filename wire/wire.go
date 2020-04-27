//+build wireinject

package wire

import (
	"os"

	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitCountersService(db *gorm.DB) services.CountersService {
	wire.Build(services.NewCountersService)

	return services.CountersService{}
}

func InitListsService(db *gorm.DB) services.ListsService {
	if inTestingMode() {
		return initMockedListsService()
	} else {
		return initDefaultListsService(db)
	}
}

func initDefaultListsService(db *gorm.DB) services.ListsService {
	wire.Build(ListsServiceSet)
	return nil
}

func initMockedListsService() services.ListsService {
	wire.Build(MockedListsServiceSet)
	return nil
}

func InitAuthService() services.AuthService {
	if inTestingMode() {
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

func InitUsersService(db *gorm.DB) services.UsersService {
	if inTestingMode() {
		return initMockedUsersService()
	} else {
		return initDefaultUsersService(db)
	}
}

func initDefaultUsersService(db *gorm.DB) services.UsersService {
	wire.Build(CryptoHelperSet, UsersServiceSet)
	return nil
}

func initMockedUsersService() services.UsersService {
	wire.Build(MockedUsersServiceSet)
	return nil
}

func InitConfigurationService() services.ConfigurationService {
	wire.Build(ConfigurationServiceSet)
	return nil
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
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

var UsersServiceSet = wire.NewSet(
	services.NewDefaultUsersService,
	wire.Bind(new(services.UsersService), new(*services.DefaultUsersService)))

var MockedUsersServiceSet = wire.NewSet(
	services.NewMockedUsersService,
	wire.Bind(new(services.UsersService), new(*services.MockedUsersService)))

var ListsServiceSet = wire.NewSet(
	services.NewDefaultListsService,
	wire.Bind(new(services.ListsService), new(*services.DefaultListsService)))

var MockedListsServiceSet = wire.NewSet(
	services.NewMockedListsService,
	wire.Bind(new(services.ListsService), new(*services.MockedListsService)))
