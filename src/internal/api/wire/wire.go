//+build wireinject

package wire

import (
	"os"

	"github.com/AngelVlc/todos/internal/api/handlers"
	"github.com/AngelVlc/todos/internal/api/repositories"
	"github.com/AngelVlc/todos/internal/api/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitLogMiddleware() handlers.LogMiddleware {
	wire.Build(handlers.NewLogMiddleware)
	return handlers.LogMiddleware{}
}

func InitAuthMiddleware(db *gorm.DB) handlers.AuthMiddleware {
	if inTestingMode() {
		return initMockedAuthMiddleware()
	} else {
		return initDefaultAuthMiddleware(db)
	}
}

func initDefaultAuthMiddleware(db *gorm.DB) handlers.AuthMiddleware {
	wire.Build(AuthMiddlewareSet)
	return nil
}

func initMockedAuthMiddleware() handlers.AuthMiddleware {
	wire.Build(MockedAuthMiddlewareSet)
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

func InitAuthService() services.AuthService {
	if inTestingMode() {
		return initMockedAuthService()
	} else {
		return initDefaultAuthService()
	}
}

func initDefaultAuthService() services.AuthService {
	wire.Build(AuthServiceSet)
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
	wire.Build(CryptoHelperSet, UsersRepositorySet, UsersServiceSet)
	return nil
}

func initMockedUsersService() services.UsersService {
	wire.Build(MockedUsersServiceSet)
	return nil
}

func InitUsersRepository(db *gorm.DB) repositories.UsersRepository {
	if inTestingMode() {
		return initMockedUsersRepository()
	} else {
		return initDefaultUsersRepository(db)
	}
}

func initDefaultUsersRepository(db *gorm.DB) repositories.UsersRepository {
	wire.Build(UsersRepositorySet)
	return nil
}

func initMockedUsersRepository() repositories.UsersRepository {
	wire.Build(MockedUsersRepositorySet)
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
	TokenHelperSet,
	ConfigurationServiceSet,
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
	AuthServiceSet,
	handlers.NewDefaultAuthMiddleware,
	wire.Bind(new(handlers.AuthMiddleware), new(*handlers.DefaultAuthMiddleware)))

var MockedAuthMiddlewareSet = wire.NewSet(
	handlers.NewMockedAuthMiddleware,
	wire.Bind(new(handlers.AuthMiddleware), new(*handlers.MockedAuthMiddleware)))

var RequireAdminMiddlewareSet = wire.NewSet(
	handlers.NewDefaultRequireAdminMiddleware,
	wire.Bind(new(handlers.RequireAdminMiddleware), new(*handlers.DefaultRequireAdminMiddleware)))

var UsersRepositorySet = wire.NewSet(
	repositories.NewDefaultUsersRepository,
	wire.Bind(new(repositories.UsersRepository), new(*repositories.DefaultUsersRepository)))

var MockedUsersRepositorySet = wire.NewSet(
	repositories.NewMockedUsersRepository,
	wire.Bind(new(repositories.UsersRepository), new(*repositories.MockedUsersRepository)))

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
