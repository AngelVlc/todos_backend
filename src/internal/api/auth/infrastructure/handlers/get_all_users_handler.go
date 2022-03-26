package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	srv := application.NewGetAllUsersService(h.AuthRepository)
	foundUsers, err := srv.GetAllUsers(r.Context())
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	res := make([]*infrastructure.UserResponse, len(foundUsers))

	for i, v := range foundUsers {
		res[i] = &infrastructure.UserResponse{
			ID:      v.ID,
			Name:    string(v.Name),
			IsAdmin: v.IsAdmin,
		}
	}

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}