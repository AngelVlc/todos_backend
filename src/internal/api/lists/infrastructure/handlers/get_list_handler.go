package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
)

func GetListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewGetListService(h.ListsRepository)
	foundList, err := srv.GetList(r.Context(), listID, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.ListResponse{
		ID:         foundList.ID,
		Name:       string(foundList.Name),
		ItemsCount: foundList.ItemsCount,
	}

	return results.OkResult{Content: &res, StatusCode: http.StatusOK}
}
