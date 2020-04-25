package controllers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

func AddUserHandler(r *http.Request, db *gorm.DB) HandlerResult {
	var dto dtos.UserDto
	err := parseBody(r, &dto)
	if err != nil {
		return errorResult{err}
	}

	userSrv := wire.InitUsersService(db)
	id, err := userSrv.AddUser(&dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}
