//+build wireinject

package wire

import (
	"os"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/handlers"
	listsDomain "github.com/AngelVlc/todos/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	authMiddleware "github.com/AngelVlc/todos/internal/api/server/middlewares/auth"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	sharedDomain "github.com/AngelVlc/todos/internal/api/shared/domain"
	sharedInfra "github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	sharedMiddlewares "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitLogMiddleware() handlers.LogMiddleware {
	wire.Build(handlers.NewLogMiddleware)
	return handlers.LogMiddleware{}
}

func InitAuthMiddleware(db *gorm.DB) authMiddleware.AuthMiddleware {
	if inTestingMode() {
		return initFakeAuthMiddleware()
	} else {
		return initDefaultAuthMiddleware()
	}
}

func initDefaultAuthMiddleware() authMiddleware.AuthMiddleware {
	wire.Build(AuthMiddlewareSet)
	return nil
}

func initFakeAuthMiddleware() authMiddleware.AuthMiddleware {
	wire.Build(FakeAuthMiddlewareSet)
	return nil
}

func InitRequireAdminMiddleware() handlers.RequireAdminMiddleware {
	wire.Build(RequireAdminMiddlewareSet)
	return nil
}

func InitRequestCounterMiddleware(db *gorm.DB) sharedMiddlewares.RequestCounterMiddleware {
	if inTestingMode() {
		return initMockedRequestCounterMiddleware()
	} else {
		return initDefaultRequestCounterMiddleware(db)
	}
}

func initDefaultRequestCounterMiddleware(db *gorm.DB) sharedMiddlewares.RequestCounterMiddleware {
	wire.Build(RequestCounterMiddlewareSet)
	return nil
}

func initMockedRequestCounterMiddleware() sharedMiddlewares.RequestCounterMiddleware {
	wire.Build(MockedRequestCounterMiddlewareSet)
	return nil
}

func InitConfigurationService() sharedApp.ConfigurationService {
	wire.Build(ConfigurationServiceSet)
	return nil
}

func InitAuthRepository(db *gorm.DB) authDomain.AuthRepository {
	if inTestingMode() {
		return initMockedAuthRepositorySet()
	} else {
		return initMySqlAuthRepository(db)
	}
}

func initMockedAuthRepositorySet() authDomain.AuthRepository {
	wire.Build(MockedAuthRepositorySet)
	return nil
}

func initMySqlAuthRepository(db *gorm.DB) authDomain.AuthRepository {
	wire.Build(MySqlAuthRepositorySet)
	return nil
}

func InitPasswordGenerator() authDomain.PasswordGenerator {
	if inTestingMode() {
		return initMockedPasswordGenerator()
	} else {
		return initBryptPasswordGenerator()
	}
}

func initBryptPasswordGenerator() authDomain.PasswordGenerator {
	wire.Build(BcryptPasswordGeneratorSet)
	return nil
}

func initMockedPasswordGenerator() authDomain.PasswordGenerator {
	wire.Build(MockedPasswordGeneratorSet)
	return nil
}

func InitListsRepository(db *gorm.DB) listsDomain.ListsRepository {
	if inTestingMode() {
		return initMockedListsRepositorySet()
	} else {
		return initMySqlListsRepository(db)
	}
}

func initMockedListsRepositorySet() listsDomain.ListsRepository {
	wire.Build(MockedListsRepositorySet)
	return nil
}

func initMySqlListsRepository(db *gorm.DB) listsDomain.ListsRepository {
	wire.Build(MySqlListsRepositorySet)
	return nil
}

func InitCountersRepository(db *gorm.DB) sharedDomain.CountersRepository {
	if inTestingMode() {
		return initMockedCountersRepositorySet()
	} else {
		return initMySqlCountersRepository(db)
	}
}

func initMockedCountersRepositorySet() sharedDomain.CountersRepository {
	wire.Build(MockedCountersRepositorySet)
	return nil
}

func initMySqlCountersRepository(db *gorm.DB) sharedDomain.CountersRepository {
	wire.Build(MySqlCountersRepositorySet)
	return nil
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var EnvGetterSet = wire.NewSet(
	sharedApp.NewOsEnvGetter,
	wire.Bind(new(sharedApp.EnvGetter), new(*sharedApp.OsEnvGetter)))

var ConfigurationServiceSet = wire.NewSet(
	EnvGetterSet,
	sharedApp.NewDefaultConfigurationService,
	wire.Bind(new(sharedApp.ConfigurationService), new(*sharedApp.DefaultConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(
	sharedApp.NewMockedConfigurationService,
	wire.Bind(new(sharedApp.ConfigurationService), new(*sharedApp.MockedConfigurationService)))

var RequestCounterMiddlewareSet = wire.NewSet(
	MySqlCountersRepositorySet,
	sharedMiddlewares.NewDefaultRequestCounterMiddleware,
	wire.Bind(new(sharedMiddlewares.RequestCounterMiddleware), new(*sharedMiddlewares.DefaultRequestCounterMiddleware)))

var MockedRequestCounterMiddlewareSet = wire.NewSet(
	sharedMiddlewares.NewMockedRequestCounterMiddleware,
	wire.Bind(new(sharedMiddlewares.RequestCounterMiddleware), new(*sharedMiddlewares.MockedRequestCounterMiddleware)))

var AuthMiddlewareSet = wire.NewSet(
	ConfigurationServiceSet,
	authMiddleware.NewRealAuthMiddleware,
	wire.Bind(new(authMiddleware.AuthMiddleware), new(*authMiddleware.RealAuthMiddleware)))

var FakeAuthMiddlewareSet = wire.NewSet(
	authMiddleware.NewFakeAuthMiddleware,
	wire.Bind(new(authMiddleware.AuthMiddleware), new(*authMiddleware.FakeAuthMiddleware)))

var RequireAdminMiddlewareSet = wire.NewSet(
	handlers.NewDefaultRequireAdminMiddleware,
	wire.Bind(new(handlers.RequireAdminMiddleware), new(*handlers.DefaultRequireAdminMiddleware)))

var MySqlAuthRepositorySet = wire.NewSet(
	authRepository.NewMySqlAuthRepository,
	wire.Bind(new(authDomain.AuthRepository), new(*authRepository.MySqlAuthRepository)))

var MockedAuthRepositorySet = wire.NewSet(
	authRepository.NewMockedAuthRepository,
	wire.Bind(new(authDomain.AuthRepository), new(*authRepository.MockedAuthRepository)))

var BcryptPasswordGeneratorSet = wire.NewSet(
	authDomain.NewBcryptPasswordGenerator,
	wire.Bind(new(authDomain.PasswordGenerator), new(*authDomain.BcryptPasswordGenerator)))

var MockedPasswordGeneratorSet = wire.NewSet(
	authDomain.NewMockedPasswordGenerator,
	wire.Bind(new(authDomain.PasswordGenerator), new(*authDomain.MockedPasswordGenerator)))

var MySqlListsRepositorySet = wire.NewSet(
	listsRepository.NewMySqlListsRepository,
	wire.Bind(new(listsDomain.ListsRepository), new(*listsRepository.MySqlListsRepository)))

var MockedListsRepositorySet = wire.NewSet(
	listsRepository.NewMockedListsRepository,
	wire.Bind(new(listsDomain.ListsRepository), new(*listsRepository.MockedListsRepository)))

var MySqlCountersRepositorySet = wire.NewSet(
	sharedInfra.NewMySqlCountersRepository,
	wire.Bind(new(sharedDomain.CountersRepository), new(*sharedInfra.MySqlCountersRepository)))

var MockedCountersRepositorySet = wire.NewSet(
	sharedInfra.NewMockedCountersRepository,
	wire.Bind(new(sharedDomain.CountersRepository), new(*sharedInfra.MockedCountersRepository)))
