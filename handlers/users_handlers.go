package handlers

import (
	"fmt"
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

func GetUsersHandler(r *http.Request, h Handler) HandlerResult {
	res := []dtos.GetUsersResultDto{}
	err := h.usersSrv.GetUsers(&res)
	if err != nil {
		return errorResult{err}
	}
	fmt.Println(res)
	return okResult{res, http.StatusOK}
}
