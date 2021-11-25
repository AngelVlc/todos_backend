//+build wireinject

package wire

import (
	"os"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	passgen "github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	authRepository "github.com/AngelVlc/todos/internal/api/auth/infrastructure/repository"
	listsDomain "github.com/AngelVlc/todos/internal/api/lists/domain"
	listsRepository "github.com/AngelVlc/todos/internal/api/lists/infrastructure/repository"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	sharedDomain "github.com/AngelVlc/todos/internal/api/shared/domain"
	events "github.com/AngelVlc/todos/internal/api/shared/domain/events"
	sharedInfra "github.com/AngelVlc/todos/internal/api/shared/infrastructure"
	authMiddleware "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/auth"
	fakemdw "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/fake"
	logMdw "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/log"
	reqadminmdw "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqadmin"
	reqcountermdw "github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/reqcounter"
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

func InitRequestCounterMiddleware(db *gorm.DB) sharedDomain.Middleware {
	if inTestingMode() {
		return initFakeMiddleware()
	} else {
		return initRequestCounterMiddleware(db)
	}
}

func initRequestCounterMiddleware(db *gorm.DB) sharedDomain.Middleware {
	wire.Build(RequestCounterMiddlewareSet)
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

func InitCountersRepository(db *gorm.DB) sharedDomain.CountersRepository {
	if inTestingMode() {
		return initMockedCountersRepository()
	} else {
		return initMySqlCountersRepository(db)
	}
}

func initMockedCountersRepository() sharedDomain.CountersRepository {
	wire.Build(MockedCountersRepositorySet)
	return nil
}

func initMySqlCountersRepository(db *gorm.DB) sharedDomain.CountersRepository {
	wire.Build(MySqlCountersRepositorySet)
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

var RequestCounterMiddlewareSet = wire.NewSet(
	MySqlCountersRepositorySet,
	sharedApp.NewIncrementRequestsCounterService,
	reqcountermdw.NewRequestCounterMiddleware,
	wire.Bind(new(sharedDomain.Middleware), new(*reqcountermdw.RequestCounterMiddleware)))

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

var MySqlCountersRepositorySet = wire.NewSet(
	sharedInfra.NewMySqlCountersRepository,
	wire.Bind(new(sharedDomain.CountersRepository), new(*sharedInfra.MySqlCountersRepository)))

var MockedCountersRepositorySet = wire.NewSet(
	sharedInfra.NewMockedCountersRepository,
	wire.Bind(new(sharedDomain.CountersRepository), new(*sharedInfra.MockedCountersRepository)))

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
