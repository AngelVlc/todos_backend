package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

// RefreshTokenHandler is the handler for the auth/refreshtoken endpoint
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	rt, err := getRefreshTokenCookieValue(r)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewRefreshTokenService(h.AuthRepository, h.CfgSrv)
	res, err := srv.RefreshToken(rt)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	addRefreshTokenCookie(w, res.RefreshToken)

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}

func getRefreshTokenCookieValue(r *http.Request) (string, error) {
	c, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", &appErrors.BadRequestError{Msg: "Missing refresh token cookie", InternalError: err}
	}

	return c.Value, nil
}
