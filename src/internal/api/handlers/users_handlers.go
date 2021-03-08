package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/dtos"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	var dto dtos.UserDto
	err := helpers.ParseBody(r, &dto)
	if err != nil {
		return results.ErrorResult{err}
	}

	id, err := h.UsersSrv.AddUser(&dto)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{id, http.StatusCreated}
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	foundUserLists, err := h.ListsSrv.GetUserLists(userID)
	if err != nil {
		return results.ErrorResult{err}
	}

	if len(foundUserLists) > 0 {
		return results.ErrorResult{&appErrors.BadRequestError{Msg: "The user has lists", InternalError: nil}}
	}

	err = h.UsersSrv.RemoveUser(userID)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{nil, http.StatusNoContent}
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	var dto dtos.UserDto
	err := helpers.ParseBody(r, &dto)
	if err != nil {
		return results.ErrorResult{err}
	}

	err = h.UsersSrv.UpdateUser(userID, &dto)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{nil, http.StatusNoContent}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.ParseInt32UrlVar(r, "id")

	u, err := h.UsersSrv.FindUserByID(userID)
	if err != nil {
		return results.ErrorResult{err}
	}

	result := dtos.UserResponseDto{
		Name:    u.Name,
		IsAdmin: u.IsAdmin,
		ID:      u.ID,
	}

	return results.OkResult{result, http.StatusOK}
}
