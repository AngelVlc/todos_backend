package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func DeleteRefreshTokensHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	input, _ := h.RequestInput.(*[]int32)

	srv := application.NewDeleteRefreshTokensService(h.AuthRepository)
	err := srv.DeleteRefreshTokens(r.Context(), *input)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
