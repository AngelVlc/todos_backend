package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := h.ParseInt32UrlVar(r, "id")
	userID := h.GetUserIDFromContext(r)

	srv := application.NewGetListService(h.ListsRepository)
	foundList, err := srv.GetList(r.Context(), listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: foundList, StatusCode: http.StatusOK}
}
