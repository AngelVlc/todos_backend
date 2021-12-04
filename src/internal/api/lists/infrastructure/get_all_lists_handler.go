package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func GetAllListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	srv := application.NewGetAllListsService(h.ListsRepository)
	foundLists, err := srv.GetAllLists(r.Context(), userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := make([]ListResponse, len(foundLists))

	for i, v := range foundLists {
		res[i] = ListResponse{
			ID:         v.ID,
			Name:       string(v.Name),
			ItemsCount: v.ItemsCount,
		}
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
