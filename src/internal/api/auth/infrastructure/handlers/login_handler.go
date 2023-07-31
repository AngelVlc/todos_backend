package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

// LoginHandler is the handler for the /auth/login endpoint
func LoginHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	input, _ := h.RequestInput.(*infrastructure.LoginInput)

	srv := application.NewLoginService(h.AuthRepository, h.UsersRepository, h.CfgSrv, h.TokenSrv)
	res, err := srv.Login(r.Context(), input.UserName.String(), input.Password.String())
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
