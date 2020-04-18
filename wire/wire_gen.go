// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/AngelVlc/todos/services"
	"github.com/jinzhu/gorm"
)

// Injectors from wire.go:

func InitUsersService(db *gorm.DB) services.UsersService {
	usersService := services.NewUsersService(db)
	return usersService
}

func InitCountersService(db *gorm.DB) services.CountersService {
	countersService := services.NewCountersService(db)
	return countersService
}

func InitListsService(db *gorm.DB) services.ListsService {
	listsService := services.NewListsService(db)
	return listsService
}

func InitJwtProvider() services.JwtProvider {
	configurationService := services.NewConfigurationService()
	jwtProvider := services.NewJwtProvider(configurationService)
	return jwtProvider
}

func InitAuthService() services.AuthService {
	configurationService := services.NewConfigurationService()
	jwtProvider := services.NewJwtProvider(configurationService)
	authService := services.NewAuthService(jwtProvider)
	return authService
}
