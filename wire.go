//+build wireinject

package main

import (
	"github.com/AngelVlc/todos/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func initUsersService(db *gorm.DB) services.UsersService {
	wire.Build(services.NewUsersService)

	return services.UsersService{}
}

func initCountersService(db *gorm.DB) services.CountersService {
	wire.Build(services.NewCountersService)

	return services.CountersService{}
}
