//+build wireinject

package wire

import (
	"os"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	"github.com/AngelVlc/todos/internal/api/handlers"
	"github.com/AngelVlc/todos/internal/api/repositories"
	authMiddleware "github.com/AngelVlc/todos/internal/api/server/middlewares/auth"
	"github.com/AngelVlc/todos/internal/api/services"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
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

func InitRequestCounterMiddleware(db *gorm.DB) handlers.RequestCounterMiddleware {
	if inTestingMode() {
		return initMockedRequestCounterMiddleware()
	} else {
		return initDefaultRequestCounterMiddleware(db)
	}
}

func initDefaultRequestCounterMiddleware(db *gorm.DB) handlers.RequestCounterMiddleware {
	wire.Build(RequestCounterMiddlewareSet)
	return nil
}

func initMockedRequestCounterMiddleware() handlers.RequestCounterMiddleware {
	wire.Build(MockedRequestCounterMiddlewareSet)
	return nil
}

func InitCountersService(db *gorm.DB) services.CountersService {
	if inTestingMode() {
		return initMockedCountersService()
	} else {
		return initDefaultCountersService(db)
	}
}

func initDefaultCountersService(db *gorm.DB) services.CountersService {
	wire.Build(CountersServiceSet)
	return nil
}

func initMockedCountersService() services.CountersService {
	wire.Build(MockedCountersServiceSet)
	return nil
}

func InitListsService(db *gorm.DB) services.ListsService {
	if inTestingMode() {
		return initMockedListsService()
	} else {
		return initDefaultListsService(db)
	}
}

func initDefaultListsService(db *gorm.DB) services.ListsService {
	wire.Build(ListsRepositorySet, ListsServiceSet)
	return nil
}

func initMockedListsService() services.ListsService {
	wire.Build(MockedListsServiceSet)
	return nil
}

func InitListsRepository(db *gorm.DB) repositories.ListsRepository {
	if inTestingMode() {
		return initMockedListsRepository()
	} else {
		return initDefaultListsRepository(db)
	}
}

func initDefaultListsRepository(db *gorm.DB) repositories.ListsRepository {
	wire.Build(ListsRepositorySet)
	return nil
}

func initMockedListsRepository() repositories.ListsRepository {
	wire.Build(MockedListsRepositorySet)
	return nil
}

func InitListItemsRepository(db *gorm.DB) repositories.ListItemsRepository {
	if inTestingMode() {
		return initMockedListItemsRepository()
	} else {
		return initDefaultListItemsRepository(db)
	}
}

func initDefaultListItemsRepository(db *gorm.DB) repositories.ListItemsRepository {
	wire.Build(ListItemsRepositorySet)
	return nil
}

func initMockedListItemsRepository() repositories.ListItemsRepository {
	wire.Build(MockedListItemsRepositorySet)
	return nil
}

func InitListItemsService(db *gorm.DB) services.ListItemsService {
	if inTestingMode() {
		return initMockedListItemsService()
	} else {
		return initDefaultListItemsService(db)
	}
}

func initDefaultListItemsService(db *gorm.DB) services.ListItemsService {
	wire.Build(ListItemsRepositorySet, ListsRepositorySet, ListItemsServiceSet)
	return nil
}

func initMockedListItemsService() services.ListItemsService {
	wire.Build(MockedListItemsServiceSet)
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

var ListsServiceSet = wire.NewSet(
	services.NewDefaultListsService,
	wire.Bind(new(services.ListsService), new(*services.DefaultListsService)))

var MockedListsServiceSet = wire.NewSet(
	services.NewMockedListsService,
	wire.Bind(new(services.ListsService), new(*services.MockedListsService)))

var CountersServiceSet = wire.NewSet(
	services.NewDefaultCountersService,
	wire.Bind(new(services.CountersService), new(*services.DefaultCountersService)))

var MockedCountersServiceSet = wire.NewSet(
	services.NewMockedCountersService,
	wire.Bind(new(services.CountersService), new(*services.MockedCountersService)))

var RequestCounterMiddlewareSet = wire.NewSet(
	CountersServiceSet,
	handlers.NewDefaultRequestCounterMiddleware,
	wire.Bind(new(handlers.RequestCounterMiddleware), new(*handlers.DefaultRequestCounterMiddleware)))

var MockedRequestCounterMiddlewareSet = wire.NewSet(
	handlers.NewMockedRequestCounterMiddleware,
	wire.Bind(new(handlers.RequestCounterMiddleware), new(*handlers.MockedRequestCounterMiddleware)))

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

var ListsRepositorySet = wire.NewSet(
	repositories.NewDefaultListsRepository,
	wire.Bind(new(repositories.ListsRepository), new(*repositories.DefaultListsRepository)))

var MockedListsRepositorySet = wire.NewSet(
	repositories.NewMockedListsRepository,
	wire.Bind(new(repositories.ListsRepository), new(*repositories.MockedListsRepository)))

var ListItemsServiceSet = wire.NewSet(
	services.NewDefaultListItemsService,
	wire.Bind(new(services.ListItemsService), new(*services.DefaultListItemsService)))

var MockedListItemsServiceSet = wire.NewSet(
	services.NewMockedListItemsService,
	wire.Bind(new(services.ListItemsService), new(*services.MockedListItemsService)))

var ListItemsRepositorySet = wire.NewSet(
	repositories.NewDefaultListItemsRepository,
	wire.Bind(new(repositories.ListItemsRepository), new(*repositories.DefaultListItemsRepository)))

var MockedListItemsRepositorySet = wire.NewSet(
	repositories.NewMockedListItemsRepository,
	wire.Bind(new(repositories.ListItemsRepository), new(*repositories.MockedListItemsRepository)))

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
