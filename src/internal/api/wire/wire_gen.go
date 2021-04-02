// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package wire

import (
	domain2 "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	"github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	domain3 "github.com/AngelVlc/todos/internal/api/lists/domain"
	repository2 "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/auth"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/fake"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/log"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqadmin"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqcounter"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"os"
)

// Injectors from wire.go:

func initFakeMiddleware() domain.Middleware {
	fakeMiddleware := fakemdw.NewFakeMiddleware()
	return fakeMiddleware
}

func initLogMiddleware() domain.Middleware {
	logMiddleware := logmdw.NewLogMiddleware()
	return logMiddleware
}

func initDefaultAuthMiddleware() authmdw.AuthMiddleware {
	osEnvGetter := application.NewOsEnvGetter()
	realConfigurationService := application.NewRealConfigurationService(osEnvGetter)
	realAuthMiddleware := authmdw.NewRealAuthMiddleware(realConfigurationService)
	return realAuthMiddleware
}

func initFakeAuthMiddleware() authmdw.AuthMiddleware {
	fakeAuthMiddleware := authmdw.NewFakeAuthMiddleware()
	return fakeAuthMiddleware
}

func InitRequireAdminMiddleware() domain.Middleware {
	requireAdminMiddleware := reqadminmdw.NewRequireAdminMiddleware()
	return requireAdminMiddleware
}

func initRequestCounterMiddleware(db *gorm.DB) domain.Middleware {
	mySqlCountersRepository := infrastructure.NewMySqlCountersRepository(db)
	requestCounterMiddleware := reqcountermdw.NewRequestCounterMiddleware(mySqlCountersRepository)
	return requestCounterMiddleware
}

func InitConfigurationService() application.ConfigurationService {
	osEnvGetter := application.NewOsEnvGetter()
	realConfigurationService := application.NewRealConfigurationService(osEnvGetter)
	return realConfigurationService
}

func initMockedAuthRepositorySet() domain2.AuthRepository {
	mockedAuthRepository := repository.NewMockedAuthRepository()
	return mockedAuthRepository
}

func initMySqlAuthRepository(db *gorm.DB) domain2.AuthRepository {
	mySqlAuthRepository := repository.NewMySqlAuthRepository(db)
	return mySqlAuthRepository
}

func initBryptPasswordGenerator() passgen.PasswordGenerator {
	bcryptPasswordGenerator := passgen.NewBcryptPasswordGenerator()
	return bcryptPasswordGenerator
}

func initMockedPasswordGenerator() passgen.PasswordGenerator {
	mockedPasswordGenerator := passgen.NewMockedPasswordGenerator()
	return mockedPasswordGenerator
}

func initMockedListsRepositorySet() domain3.ListsRepository {
	mockedListsRepository := repository2.NewMockedListsRepository()
	return mockedListsRepository
}

func initMySqlListsRepository(db *gorm.DB) domain3.ListsRepository {
	mySqlListsRepository := repository2.NewMySqlListsRepository(db)
	return mySqlListsRepository
}

func initMockedCountersRepositorySet() domain.CountersRepository {
	mockedCountersRepository := infrastructure.NewMockedCountersRepository()
	return mockedCountersRepository
}

func initMySqlCountersRepository(db *gorm.DB) domain.CountersRepository {
	mySqlCountersRepository := infrastructure.NewMySqlCountersRepository(db)
	return mySqlCountersRepository
}

// wire.go:

func InitLogMiddleware() domain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initLogMiddleware()
	}
}

func InitAuthMiddleware(db *gorm.DB) authmdw.AuthMiddleware {
	if inTestingMode() {
		return initFakeAuthMiddleware()
	} else {
		return initDefaultAuthMiddleware()
	}
}

func InitRequestCounterMiddleware(db *gorm.DB) domain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initRequestCounterMiddleware(db)
	}
}

func InitAuthRepository(db *gorm.DB) domain2.AuthRepository {
	if inTestingMode() {
		return initMockedAuthRepositorySet()
	} else {
		return initMySqlAuthRepository(db)
	}
}

func InitPasswordGenerator() passgen.PasswordGenerator {
	if inTestingMode() {
		return initMockedPasswordGenerator()
	} else {
		return initBryptPasswordGenerator()
	}
}

func InitListsRepository(db *gorm.DB) domain3.ListsRepository {
	if inTestingMode() {
		return initMockedListsRepositorySet()
	} else {
		return initMySqlListsRepository(db)
	}
}

func InitCountersRepository(db *gorm.DB) domain.CountersRepository {
	if inTestingMode() {
		return initMockedCountersRepositorySet()
	} else {
		return initMySqlCountersRepository(db)
	}
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var EnvGetterSet = wire.NewSet(application.NewOsEnvGetter, wire.Bind(new(application.EnvGetter), new(*application.OsEnvGetter)))

var RealConfigurationServiceSet = wire.NewSet(
	EnvGetterSet, application.NewRealConfigurationService, wire.Bind(new(application.ConfigurationService), new(*application.RealConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(application.NewMockedConfigurationService, wire.Bind(new(application.ConfigurationService), new(*application.MockedConfigurationService)))

var FakeMiddlewareSet = wire.NewSet(fakemdw.NewFakeMiddleware, wire.Bind(new(domain.Middleware), new(*fakemdw.FakeMiddleware)))

var RequestCounterMiddlewareSet = wire.NewSet(
	MySqlCountersRepositorySet, reqcountermdw.NewRequestCounterMiddleware, wire.Bind(new(domain.Middleware), new(*reqcountermdw.RequestCounterMiddleware)))

var LogMiddlewareSet = wire.NewSet(logmdw.NewLogMiddleware, wire.Bind(new(domain.Middleware), new(*logmdw.LogMiddleware)))

var AuthMiddlewareSet = wire.NewSet(
	RealConfigurationServiceSet, authmdw.NewRealAuthMiddleware, wire.Bind(new(authmdw.AuthMiddleware), new(*authmdw.RealAuthMiddleware)))

var FakeAuthMiddlewareSet = wire.NewSet(authmdw.NewFakeAuthMiddleware, wire.Bind(new(authmdw.AuthMiddleware), new(*authmdw.FakeAuthMiddleware)))

var RequireAdminMiddlewareSet = wire.NewSet(reqadminmdw.NewRequireAdminMiddleware, wire.Bind(new(domain.Middleware), new(*reqadminmdw.RequireAdminMiddleware)))

var MySqlAuthRepositorySet = wire.NewSet(repository.NewMySqlAuthRepository, wire.Bind(new(domain2.AuthRepository), new(*repository.MySqlAuthRepository)))

var MockedAuthRepositorySet = wire.NewSet(repository.NewMockedAuthRepository, wire.Bind(new(domain2.AuthRepository), new(*repository.MockedAuthRepository)))

var BcryptPasswordGeneratorSet = wire.NewSet(passgen.NewBcryptPasswordGenerator, wire.Bind(new(passgen.PasswordGenerator), new(*passgen.BcryptPasswordGenerator)))

var MockedPasswordGeneratorSet = wire.NewSet(passgen.NewMockedPasswordGenerator, wire.Bind(new(passgen.PasswordGenerator), new(*passgen.MockedPasswordGenerator)))

var MySqlListsRepositorySet = wire.NewSet(repository2.NewMySqlListsRepository, wire.Bind(new(domain3.ListsRepository), new(*repository2.MySqlListsRepository)))

var MockedListsRepositorySet = wire.NewSet(repository2.NewMockedListsRepository, wire.Bind(new(domain3.ListsRepository), new(*repository2.MockedListsRepository)))

var MySqlCountersRepositorySet = wire.NewSet(infrastructure.NewMySqlCountersRepository, wire.Bind(new(domain.CountersRepository), new(*infrastructure.MySqlCountersRepository)))

var MockedCountersRepositorySet = wire.NewSet(infrastructure.NewMockedCountersRepository, wire.Bind(new(domain.CountersRepository), new(*infrastructure.MockedCountersRepository)))
