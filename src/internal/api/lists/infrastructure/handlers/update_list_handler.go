package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func UpdateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)

	updateReq, _ := h.RequestInput.(*domain.UpdateListInput)

	srv := application.NewUpdateListService(h.ListsRepository)
	list, err := srv.UpdateList(r.Context(), listID, updateReq.Name, userID, updateReq.IDsByPosition)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.ListResponse{
		ID:         list.ID,
		Name:       list.Name.String(),
		ItemsCount: list.ItemsCount,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
