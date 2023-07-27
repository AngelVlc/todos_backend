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

func CreateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	createReq, _ := h.RequestInput.(*domain.CreateListInput)

	srv := application.NewCreateListService(h.ListsRepository)
	newList, err := srv.CreateList(r.Context(), createReq.Name, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := infrastructure.ListResponse{
		ID:         newList.ID,
		Name:       newList.Name.String(),
		ItemsCount: newList.ItemsCount,
	}

	return results.OkResult{Content: res, StatusCode: http.StatusCreated}
}
