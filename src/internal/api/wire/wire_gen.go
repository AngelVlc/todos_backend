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
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/auth"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/fake"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/log"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqadmin"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqcounter"
	"github.com/google/wire"
	"gorm.io/gorm"
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
	realConfigurationService := application.NewRealConfigurationService()
	realTokenService := domain2.NewRealTokenService(realConfigurationService)
	realAuthMiddleware := authmdw.NewRealAuthMiddleware(realTokenService)
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

func initRequestIdMiddleware(db *gorm.DB) domain.Middleware {
	requestIdMiddleware := reqcountermdw.NewRequestIdMiddleware()
	return requestIdMiddleware
}

func InitConfigurationService() application.ConfigurationService {
	realConfigurationService := application.NewRealConfigurationService()
	return realConfigurationService
}

func initMockedAuthRepository() domain2.AuthRepository {
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

func initMockedListsRepository() domain3.ListsRepository {
	mockedListsRepository := repository2.NewMockedListsRepository()
	return mockedListsRepository
}

func initMySqlListsRepository(db *gorm.DB) domain3.ListsRepository {
	mySqlListsRepository := repository2.NewMySqlListsRepository(db)
	return mySqlListsRepository
}

func initMockedTokenService() domain2.TokenService {
	mockedTokenService := domain2.NewMockedTokenService()
	return mockedTokenService
}

func initRealTokenService() domain2.TokenService {
	realConfigurationService := application.NewRealConfigurationService()
	realTokenService := domain2.NewRealTokenService(realConfigurationService)
	return realTokenService
}

func initMockedEventBus() events.EventBus {
	mockedEventBus := events.NewMockedEventBus()
	return mockedEventBus
}

func initRealEventBus(subscribers map[string]events.DataChannelSlice) events.EventBus {
	realEventBus := events.NewRealEventBus(subscribers)
	return realEventBus
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

func InitRequestIdMiddleware(db *gorm.DB) domain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initRequestIdMiddleware(db)
	}
}

func InitAuthRepository(db *gorm.DB) domain2.AuthRepository {
	if inTestingMode() {
		return initMockedAuthRepository()
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
		return initMockedListsRepository()
	} else {
		return initMySqlListsRepository(db)
	}
}

func InitTokenService() domain2.TokenService {
	if inTestingMode() {
		return initMockedTokenService()
	} else {
		return initRealTokenService()
	}
}

func InitEventBus(subscribers map[string]events.DataChannelSlice) events.EventBus {
	if inTestingMode() {
		return initMockedEventBus()
	} else {
		return initRealEventBus(subscribers)
	}
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var RealConfigurationServiceSet = wire.NewSet(application.NewRealConfigurationService, wire.Bind(new(application.ConfigurationService), new(*application.RealConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(application.NewMockedConfigurationService, wire.Bind(new(application.ConfigurationService), new(*application.MockedConfigurationService)))

var FakeMiddlewareSet = wire.NewSet(fakemdw.NewFakeMiddleware, wire.Bind(new(domain.Middleware), new(*fakemdw.FakeMiddleware)))

var RequestIdMiddlewareSet = wire.NewSet(reqcountermdw.NewRequestIdMiddleware, wire.Bind(new(domain.Middleware), new(*reqcountermdw.RequestIdMiddleware)))

var LogMiddlewareSet = wire.NewSet(logmdw.NewLogMiddleware, wire.Bind(new(domain.Middleware), new(*logmdw.LogMiddleware)))

var AuthMiddlewareSet = wire.NewSet(
	RealTokenServiceSet, authmdw.NewRealAuthMiddleware, wire.Bind(new(authmdw.AuthMiddleware), new(*authmdw.RealAuthMiddleware)))

var FakeAuthMiddlewareSet = wire.NewSet(authmdw.NewFakeAuthMiddleware, wire.Bind(new(authmdw.AuthMiddleware), new(*authmdw.FakeAuthMiddleware)))

var RequireAdminMiddlewareSet = wire.NewSet(reqadminmdw.NewRequireAdminMiddleware, wire.Bind(new(domain.Middleware), new(*reqadminmdw.RequireAdminMiddleware)))

var MySqlAuthRepositorySet = wire.NewSet(repository.NewMySqlAuthRepository, wire.Bind(new(domain2.AuthRepository), new(*repository.MySqlAuthRepository)))

var MockedAuthRepositorySet = wire.NewSet(repository.NewMockedAuthRepository, wire.Bind(new(domain2.AuthRepository), new(*repository.MockedAuthRepository)))

var BcryptPasswordGeneratorSet = wire.NewSet(passgen.NewBcryptPasswordGenerator, wire.Bind(new(passgen.PasswordGenerator), new(*passgen.BcryptPasswordGenerator)))

var MockedPasswordGeneratorSet = wire.NewSet(passgen.NewMockedPasswordGenerator, wire.Bind(new(passgen.PasswordGenerator), new(*passgen.MockedPasswordGenerator)))

var MySqlListsRepositorySet = wire.NewSet(repository2.NewMySqlListsRepository, wire.Bind(new(domain3.ListsRepository), new(*repository2.MySqlListsRepository)))

var MockedListsRepositorySet = wire.NewSet(repository2.NewMockedListsRepository, wire.Bind(new(domain3.ListsRepository), new(*repository2.MockedListsRepository)))

var RealTokenServiceSet = wire.NewSet(
	RealConfigurationServiceSet, domain2.NewRealTokenService, wire.Bind(new(domain2.TokenService), new(*domain2.RealTokenService)))

var MockedTokenServiceSet = wire.NewSet(domain2.NewMockedTokenService, wire.Bind(new(domain2.TokenService), new(*domain2.MockedTokenService)))

var RealEventBusSet = wire.NewSet(events.NewRealEventBus, wire.Bind(new(events.EventBus), new(*events.RealEventBus)))

var MockedEventBusSet = wire.NewSet(events.NewMockedEventBus, wire.Bind(new(events.EventBus), new(*events.MockedEventBus)))
