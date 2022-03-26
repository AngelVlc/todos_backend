package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/internal/api/auth/application"
	appErrors "github.com/AngelVlc/todos_backend/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/results"
)

// RefreshTokenHandler is the handler for the auth/refreshtoken endpoint
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	rt, err := getRefreshTokenCookieValue(r)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewRefreshTokenService(h.AuthRepository, h.CfgSrv, h.TokenSrv)
	newToken, err := srv.RefreshToken(r.Context(), rt)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	addTokenCookie(w, newToken)

	return results.OkResult{Content: nil, StatusCode: http.StatusOK}
}

func getRefreshTokenCookieValue(r *http.Request) (string, error) {
	c, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", &appErrors.BadRequestError{Msg: "Missing refresh token cookie", InternalError: err}
	}

	return c.Value, nil
}
