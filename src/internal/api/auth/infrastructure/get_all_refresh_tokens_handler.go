package infrastructure

import (
	"net/http"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type RefreshTokenResponse struct {
	ID             int32     `json:"id"`
	UserID         int32     `json:"userId"`
	ExpirationDate time.Time `json:"expirationDate"`
}

func GetAllRefreshTokensHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	srv := application.NewGetAllRefreshTokensService(h.AuthRepository)
	found, err := srv.GetAllRefreshTokens()
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := make([]RefreshTokenResponse, len(found))

	for i, v := range found {
		res[i] = RefreshTokenResponse{
			ID:             v.ID,
			UserID:         v.UserID,
			ExpirationDate: v.ExpirationDate,
		}
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}
