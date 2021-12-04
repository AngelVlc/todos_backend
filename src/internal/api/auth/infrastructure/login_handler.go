package infrastructure

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/auth/application"
	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

type loginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// LoginHandler is the handler for the /auth/login endpoint
func LoginHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	loginReq := loginRequest{}
	err := h.ParseBody(r, &loginReq)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	userName, err := domain.NewUserName(loginReq.UserName)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewUserPassword(loginReq.Password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewLoginService(h.AuthRepository, h.CfgSrv, h.TokenSrv)
	res, err := srv.Login(r.Context(), userName, password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	addTokenCookie(w, res.Token)
	addRefreshTokenCookie(w, res.RefreshToken)

	res.Token = ""
	res.RefreshToken = ""

	return results.OkResult{Content: res, StatusCode: http.StatusOK}
}

func addTokenCookie(w http.ResponseWriter, token string) {
	rfCookie := http.Cookie{
		Name:     tokenCookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, &rfCookie)
}

func addRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	rfCookie := http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/auth",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, &rfCookie)
}
