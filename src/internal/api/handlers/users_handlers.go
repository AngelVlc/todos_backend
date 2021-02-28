package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/dtos"
	appErrors "github.com/AngelVlc/todos/internal/api/errors"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
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

func GetUsersHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	res, err := h.usersSrv.GetUsers()
	if err != nil {
		return errorResult{err}
	}
	return okResult{res, http.StatusOK}
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := parseInt32UrlVar(r, "id")

	foundUserLists, err := h.listsSrv.GetUserLists(userID)
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

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := parseInt32UrlVar(r, "id")

	var dto dtos.UserDto
	err := parseBody(r, &dto)
	if err != nil {
		return errorResult{err}
	}

	err = h.usersSrv.UpdateUser(userID, &dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := parseInt32UrlVar(r, "id")

	u, err := h.usersSrv.FindUserByID(userID)
	if err != nil {
		return errorResult{err}
	}

	result := dtos.UserResponseDto{
		Name:    u.Name,
		IsAdmin: u.IsAdmin,
		ID:      u.ID,
	}

	return okResult{result, http.StatusOK}
}
