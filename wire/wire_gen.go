// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"os"
)

// Injectors from wire.go:

func initDefaultCountersService(db *gorm.DB) services.CountersService {
	defaultCountersService := services.NewDefaultCountersService(db)
	return defaultCountersService
}

func initMockedCountersService() services.CountersService {
	mockedCountersService := services.NewMockedCountersService()
	return mockedCountersService
}

func initDefaultListsService(db *gorm.DB) services.ListsService {
	defaultListsService := services.NewDefaultListsService(db)
	return defaultListsService
}

func initMockedListsService() services.ListsService {
	mockedListsService := services.NewMockedListsService()
	return mockedListsService
}

func initDefaultAuthService() services.AuthService {
	jwtTokenHelper := services.NewJwtTokenHelper()
	osEnvGetter := services.NewOsEnvGetter()
	defaultConfigurationService := services.NewDefaultConfigurationService(osEnvGetter)
	defaultAuthService := services.NewDefaultAuthService(jwtTokenHelper, defaultConfigurationService)
	return defaultAuthService
}

func initMockedAuthService() services.AuthService {
	mockedAuthService := services.NewMockedAuthService()
	return mockedAuthService
}

func initDefaultUsersService(db *gorm.DB) services.UsersService {
	bcryptHelper := services.NewBcryptHelper()
	defaultUsersService := services.NewDefaultUsersService(bcryptHelper, db)
	return defaultUsersService
}

func initMockedUsersService() services.UsersService {
	mockedUsersService := services.NewMockedUsersService()
	return mockedUsersService
}

func InitConfigurationService() services.ConfigurationService {
	osEnvGetter := services.NewOsEnvGetter()
	defaultConfigurationService := services.NewDefaultConfigurationService(osEnvGetter)
	return defaultConfigurationService
}

// wire.go:

func InitCountersService(db *gorm.DB) services.CountersService {
	if inTestingMode() {
		return initMockedCountersService()
	} else {
		return initDefaultCountersService(db)
	}
}

func InitListsService(db *gorm.DB) services.ListsService {
	if inTestingMode() {
		return initMockedListsService()
	} else {
		return initDefaultListsService(db)
	}
}

func InitAuthService() services.AuthService {
	if inTestingMode() {
		return initMockedAuthService()
	} else {
		return initDefaultAuthService()
	}
}

func InitUsersService(db *gorm.DB) services.UsersService {
	if inTestingMode() {
		return initMockedUsersService()
	} else {
		return initDefaultUsersService(db)
	}
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var EnvGetterSet = wire.NewSet(services.NewOsEnvGetter, wire.Bind(new(services.EnvGetter), new(*services.OsEnvGetter)))

var TokenHelperSet = wire.NewSet(services.NewJwtTokenHelper, wire.Bind(new(services.TokenHelper), new(*services.JwtTokenHelper)))

var MockedTokenHelperSet = wire.NewSet(services.NewMockedTokenHelper, wire.Bind(new(services.TokenHelper), new(*services.MockedTokenHelper)))

var CryptoHelperSet = wire.NewSet(services.NewBcryptHelper, wire.Bind(new(services.CryptoHelper), new(*services.BcryptHelper)))

var ConfigurationServiceSet = wire.NewSet(
	EnvGetterSet, services.NewDefaultConfigurationService, wire.Bind(new(services.ConfigurationService), new(*services.DefaultConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(services.NewMockedConfigurationService, wire.Bind(new(services.ConfigurationService), new(*services.MockedConfigurationService)))

var AuthServiceSet = wire.NewSet(services.NewDefaultAuthService, wire.Bind(new(services.AuthService), new(*services.DefaultAuthService)))

var MockedAuthServiceSet = wire.NewSet(services.NewMockedAuthService, wire.Bind(new(services.AuthService), new(*services.MockedAuthService)))

var UsersServiceSet = wire.NewSet(services.NewDefaultUsersService, wire.Bind(new(services.UsersService), new(*services.DefaultUsersService)))

var MockedUsersServiceSet = wire.NewSet(services.NewMockedUsersService, wire.Bind(new(services.UsersService), new(*services.MockedUsersService)))

var ListsServiceSet = wire.NewSet(services.NewDefaultListsService, wire.Bind(new(services.ListsService), new(*services.DefaultListsService)))

var MockedListsServiceSet = wire.NewSet(services.NewMockedListsService, wire.Bind(new(services.ListsService), new(*services.MockedListsService)))

var CountersServiceSet = wire.NewSet(services.NewDefaultCountersService, wire.Bind(new(services.CountersService), new(*services.DefaultCountersService)))

var MockedCountersServiceSet = wire.NewSet(services.NewMockedCountersService, wire.Bind(new(services.CountersService), new(*services.MockedCountersService)))
