package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type updateListRequest struct {
	Name string `json:"name"`
}

func UpdateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	listID := helpers.ParseInt32UrlVar(r, "id")
	userID := helpers.GetUserIDFromContext(r)

	updateReq := updateListRequest{}
	err := h.ParseBody(r, &updateReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listName, err := domain.NewListName(updateReq.Name)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewUpdateListService(h.ListsRepository)
	list, err := srv.UpdateList(listID, listName, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := ListResponse{
		ID:   list.ID,
		Name: string(list.Name),
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
