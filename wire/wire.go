//+build wireinject

package wire

import (
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitUsersService(db *gorm.DB) services.UsersService {
	wire.Build(services.NewUsersService)

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

func InitJwtProvider() services.JwtProvider {
	wire.Build(services.NewConfigurationService, services.NewJwtProvider)

	return services.JwtProvider{}
}

func InitAuthService() services.AuthService {
	wire.Build(services.NewJwtProvider, services.NewConfigurationService, services.NewAuthService)

	return services.AuthService{}
}
