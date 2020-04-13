//+build wireinject

package main

import (
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func initUsersService(db *gorm.DB, config *services.ConfigurationService) services.UsersService {
	wire.Build(services.NewUsersService)

	return services.UsersService{}
}
