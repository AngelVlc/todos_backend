package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
)

func AddUserHandler(r *http.Request, h Handler) HandlerResult {
	var dto dtos.UserDto
	err := parseBody(r, &dto)
	if err != nil {
		return errorResult{err}
	}

	id, err := h.usersSrv.AddUser(&dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}
