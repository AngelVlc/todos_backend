package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/lists/application"
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type createListRequest struct {
	Name string `json:"name"`
}

func CreateListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	createReq := createListRequest{}
	err := h.ParseBody(r, &createReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	listName, err := domain.NewListName(createReq.Name)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewCreateListService(h.ListsRepository)
	newList, err := srv.CreateList(listName, userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := ListResponse{
		ID:   newList.ID,
		Name: string(newList.Name),
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
