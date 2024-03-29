//go:build wireinject
// +build wireinject

package wire

import (
	"os"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	passgen "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	authRepository "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/repository"
	listsDomain "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/repository"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	events "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	authMiddleware "github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/auth"
	fakemdw "github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/fake"
	logMdw "github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/log"
	reqadminmdw "github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/reqadmin"
	reqid "github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/reqid"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	algoliaSearch "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initFakeMiddleware() sharedDomain.Middleware {
	wire.Build(FakeMiddlewareSet)
	return nil
}

func InitLogMiddleware() sharedDomain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initLogMiddleware()
	}
}

func initLogMiddleware() sharedDomain.Middleware {
	wire.Build(LogMiddlewareSet)
	return nil
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

func InitRequireAdminMiddleware() sharedDomain.Middleware {
	wire.Build(RequireAdminMiddlewareSet)
	return nil
}

func InitRequestIdMiddleware(db *gorm.DB) sharedDomain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initRequestIdMiddleware(db)
	}
}

func initRequestIdMiddleware(db *gorm.DB) sharedDomain.Middleware {
	wire.Build(RequestIdMiddlewareSet)
	return nil
}

func InitConfigurationService() sharedApp.ConfigurationService {
	wire.Build(RealConfigurationServiceSet)
	return nil
}

func InitAuthRepository(db *gorm.DB) authDomain.AuthRepository {
	if inTestingMode() {
		return initMockedAuthRepository()
	} else {
		return initMySqlAuthRepository(db)
	}
}

func initMockedAuthRepository() authDomain.AuthRepository {
	wire.Build(MockedAuthRepositorySet)
	return nil
}

func initMySqlAuthRepository(db *gorm.DB) authDomain.AuthRepository {
	wire.Build(MySqlAuthRepositorySet)
	return nil
}

func InitUsersRepository(db *gorm.DB) authDomain.UsersRepository {
	if inTestingMode() {
		return initMockedUsersRepository()
	} else {
		return initMySqlUsersRepository(db)
	}
}

func initMockedUsersRepository() authDomain.UsersRepository {
	wire.Build(MockedUsersRepositorySet)
	return nil
}

func initMySqlUsersRepository(db *gorm.DB) authDomain.UsersRepository {
	wire.Build(MySqlUsersRepositorySet)
	return nil
}

func InitPasswordGenerator() passgen.PasswordGenerator {
	if inTestingMode() {
		return initMockedPasswordGenerator()
	} else {
		return initBryptPasswordGenerator()
	}
}

func initBryptPasswordGenerator() passgen.PasswordGenerator {
	wire.Build(BcryptPasswordGeneratorSet)
	return nil
}

func initMockedPasswordGenerator() passgen.PasswordGenerator {
	wire.Build(MockedPasswordGeneratorSet)
	return nil
}

func InitListsRepository(db *gorm.DB) listsDomain.ListsRepository {
	if inTestingMode() {
		return initMockedListsRepository()
	} else {
		return initMySqlListsRepository(db)
	}
}

func initMockedListsRepository() listsDomain.ListsRepository {
	wire.Build(MockedListsRepositorySet)
	return nil
}

func initMySqlListsRepository(db *gorm.DB) listsDomain.ListsRepository {
	wire.Build(MySqlListsRepositorySet)
	return nil
}

func InitTokenService() authDomain.TokenService {
	if inTestingMode() {
		return initMockedTokenService()
	} else {
		return initRealTokenService()
	}
}

func initMockedTokenService() authDomain.TokenService {
	wire.Build(MockedTokenServiceSet)
	return nil
}

func initRealTokenService() authDomain.TokenService {
	wire.Build(RealTokenServiceSet)
	return nil
}

func InitEventBus(subscribers map[string]events.DataChannelSlice) events.EventBus {
	if inTestingMode() {
		return initMockedEventBus()
	} else {
		return initRealEventBus(subscribers)
	}
}

func initMockedEventBus() events.EventBus {
	wire.Build(MockedEventBusSet)
	return nil
}

func initRealEventBus(subscribers map[string]events.DataChannelSlice) events.EventBus {
	wire.Build(RealEventBusSet)
	return nil
}

func InitSearchIndexClient(indexName string, settings algoliaSearch.Settings) search.SearchIndexClient {
	if inTestingMode() {
		return initMockedSearchIndexClient()
	} else {
		return initAlgoliaIndexClient(indexName, settings)
	}
}

func initMockedSearchIndexClient() search.SearchIndexClient {
	wire.Build(MockedSearchIndexClientSet)
	return nil
}

func initAlgoliaIndexClient(indexName string, settings algoliaSearch.Settings) search.SearchIndexClient {
	wire.Build(AlgoliaIndexClientSet)
	return nil
}

func InitCategoriesRepository(db *gorm.DB) listsDomain.CategoriesRepository {
	if inTestingMode() {
		return initMockedCategoriesRepository()
	} else {
		return initMySqlCategoriesRepository(db)
	}
}

func initMockedCategoriesRepository() listsDomain.CategoriesRepository {
	wire.Build(MockedCategoriesRepositorySet)
	return nil
}

func initMySqlCategoriesRepository(db *gorm.DB) listsDomain.CategoriesRepository {
	wire.Build(MySqlCategoriesRepositorySet)
	return nil
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var RealConfigurationServiceSet = wire.NewSet(
	sharedApp.NewRealConfigurationService,
	wire.Bind(new(sharedApp.ConfigurationService), new(*sharedApp.RealConfigurationService)))

var MockedConfigurationServiceSet = wire.NewSet(
	sharedApp.NewMockedConfigurationService,
	wire.Bind(new(sharedApp.ConfigurationService), new(*sharedApp.MockedConfigurationService)))

var FakeMiddlewareSet = wire.NewSet(
	fakemdw.NewFakeMiddleware,
	wire.Bind(new(sharedDomain.Middleware), new(*fakemdw.FakeMiddleware)))

var RequestIdMiddlewareSet = wire.NewSet(
	reqid.NewRequestIdMiddleware,
	wire.Bind(new(sharedDomain.Middleware), new(*reqid.RequestIdMiddleware)))

var LogMiddlewareSet = wire.NewSet(
	logMdw.NewLogMiddleware,
	wire.Bind(new(sharedDomain.Middleware), new(*logMdw.LogMiddleware)))

var AuthMiddlewareSet = wire.NewSet(
	RealTokenServiceSet,
	authMiddleware.NewRealAuthMiddleware,
	wire.Bind(new(authMiddleware.AuthMiddleware), new(*authMiddleware.RealAuthMiddleware)))

var FakeAuthMiddlewareSet = wire.NewSet(
	authMiddleware.NewFakeAuthMiddleware,
	wire.Bind(new(authMiddleware.AuthMiddleware), new(*authMiddleware.FakeAuthMiddleware)))

var RequireAdminMiddlewareSet = wire.NewSet(
	reqadminmdw.NewRequireAdminMiddleware,
	wire.Bind(new(sharedDomain.Middleware), new(*reqadminmdw.RequireAdminMiddleware)))

var MySqlAuthRepositorySet = wire.NewSet(
	authRepository.NewMySqlAuthRepository,
	wire.Bind(new(authDomain.AuthRepository), new(*authRepository.MySqlAuthRepository)))

var MockedAuthRepositorySet = wire.NewSet(
	authRepository.NewMockedAuthRepository,
	wire.Bind(new(authDomain.AuthRepository), new(*authRepository.MockedAuthRepository)))

var MySqlUsersRepositorySet = wire.NewSet(
	authRepository.NewMySqlUsersRepository,
	wire.Bind(new(authDomain.UsersRepository), new(*authRepository.MySqlUsersRepository)))

var MockedUsersRepositorySet = wire.NewSet(
	authRepository.NewMockedUsersRepository,
	wire.Bind(new(authDomain.UsersRepository), new(*authRepository.MockedUsersRepository)))

var BcryptPasswordGeneratorSet = wire.NewSet(
	passgen.NewBcryptPasswordGenerator,
	wire.Bind(new(passgen.PasswordGenerator), new(*passgen.BcryptPasswordGenerator)))

var MockedPasswordGeneratorSet = wire.NewSet(
	passgen.NewMockedPasswordGenerator,
	wire.Bind(new(passgen.PasswordGenerator), new(*passgen.MockedPasswordGenerator)))

var MySqlListsRepositorySet = wire.NewSet(
	listsRepository.NewMySqlListsRepository,
	wire.Bind(new(listsDomain.ListsRepository), new(*listsRepository.MySqlListsRepository)))

var MockedListsRepositorySet = wire.NewSet(
	listsRepository.NewMockedListsRepository,
	wire.Bind(new(listsDomain.ListsRepository), new(*listsRepository.MockedListsRepository)))

var RealTokenServiceSet = wire.NewSet(
	RealConfigurationServiceSet,
	authDomain.NewRealTokenService,
	wire.Bind(new(authDomain.TokenService), new(*authDomain.RealTokenService)))

var MockedTokenServiceSet = wire.NewSet(
	authDomain.NewMockedTokenService,
	wire.Bind(new(authDomain.TokenService), new(*authDomain.MockedTokenService)))

var RealEventBusSet = wire.NewSet(
	events.NewRealEventBus,
	wire.Bind(new(events.EventBus), new(*events.RealEventBus)))

var MockedEventBusSet = wire.NewSet(
	events.NewMockedEventBus,
	wire.Bind(new(events.EventBus), new(*events.MockedEventBus)))

var AlgoliaIndexClientSet = wire.NewSet(
	RealConfigurationServiceSet,
	search.NewAlgoliaIndexClient,
	wire.Bind(new(search.SearchIndexClient), new(*search.AlgoliaIndexClient)),
)

var MockedSearchIndexClientSet = wire.NewSet(
	search.NewMockedSearchIndexClient,
	wire.Bind(new(search.SearchIndexClient), new(*search.MockedSearchIndexClient)),
)

var MySqlCategoriesRepositorySet = wire.NewSet(
	listsRepository.NewMySqlCategoriesRepository,
	wire.Bind(new(listsDomain.CategoriesRepository), new(*listsRepository.MySqlCategoriesRepository)),
)

var MockedCategoriesRepositorySet = wire.NewSet(
	listsRepository.NewMockedCategoriesRepository,
	wire.Bind(new(listsDomain.CategoriesRepository), new(*listsRepository.MockedCategoriesRepository)),
)
