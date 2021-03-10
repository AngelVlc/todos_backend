package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type loginRequest struct {
	UserName *string `json:"userName"`
	Password *string `json:"password"`
}

// LoginHandler is the handler for the /auth/login endpoint
func LoginHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	loginReq := loginRequest{}
	err := h.ParseBody(r, &loginReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	userName, err := domain.NewAuthUserName(loginReq.UserName, true)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewAuthUserPassword(loginReq.Password, true)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewLoginService(h.AuthRepository, h.CfgSrv)
	res, err := srv.Login(userName, password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	addRefreshTokenCookieKK(w, res.RefreshToken)

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}

func addRefreshTokenCookieKK(w http.ResponseWriter, refreshToken string) {
	rfCookie := http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/auth",
	}
	http.SetCookie(w, &rfCookie)
}
