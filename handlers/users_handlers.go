package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
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
	return okResult{res, http.StatusOK}
}

func DeleteUserHandler(r *http.Request, h Handler) HandlerResult {
	userID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	foundUserLists := []dtos.GetListsResultDto{}
	err = h.listsSrv.GetUserLists(userID, &foundUserLists)
	if err != nil {
		return errorResult{err}
	}

	if len(foundUserLists) > 0 {
		return errorResult{&appErrors.BadRequestError{Msg: "The user has lists", InternalError: nil}}
	}

	err = h.usersSrv.RemoveUser(userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}
