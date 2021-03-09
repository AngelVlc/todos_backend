package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

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
