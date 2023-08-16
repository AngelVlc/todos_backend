package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func GetSearchSecureKeyHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := h.GetUserIDFromContext(r)

	service := application.NewGetSearchSecureKeyService(h.SearchClient)
	key, err := service.GetSearchSecureKeyService(userID)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: key, StatusCode: http.StatusOK}
}
