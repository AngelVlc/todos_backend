package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetAllListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewGetAllListsService(h.ListsRepository)
	foundLists, err := srv.GetAllLists(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := make([]infrastructure.ListResponse, len(foundLists))

	for i, v := range foundLists {
		res[i] = infrastructure.ListResponse{
			ID:         v.ID,
			Name:       string(v.Name),
			ItemsCount: v.ItemsCount,
		}
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
