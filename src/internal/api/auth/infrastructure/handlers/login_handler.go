package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

const (
	tokenCookieName        = "token"
	refreshTokenCookieName = "refreshToken"
)

// LoginHandler is the handler for the /auth/login endpoint
func LoginHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	input, _ := h.RequestInput.(*infrastructure.LoginInput)

	srv := application.NewLoginService(h.AuthRepository, h.UsersRepository, h.CfgSrv, h.TokenSrv)
	t, rt, u, err := srv.Login(r.Context(), input.UserName, input.Password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	addTokenCookie(w, t)
	addRefreshTokenCookie(w, rt)

	res := infrastructure.UserResponse{ID: u.ID, Name: u.Name.String(), IsAdmin: u.IsAdmin}

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
