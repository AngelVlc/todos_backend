package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	"github.com/AngelVlc/todos/wire"
)

func AddUserHandler(r *http.Request, h Handler) HandlerResult {
	var dto dtos.UserDto
	err := parseBody(r, &dto)
	if err != nil {
		return errorResult{err}
	}

	userSrv := wire.InitUsersService(h.Db)
	id, err := userSrv.AddUser(&dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}
