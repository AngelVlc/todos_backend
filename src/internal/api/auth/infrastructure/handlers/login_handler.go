package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

// LoginHandler is the handler for the /auth/login endpoint
func LoginHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	loginReq, _ := h.RequestInput.(*domain.LoginInput)

	userName, err := domain.NewUserNameValueObject(loginReq.UserName)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	password, err := domain.NewUserPassword(loginReq.Password)
	if err != nil {
		return results.ErrorResult{Err: err}
	}

	srv := application.NewLoginService(h.AuthRepository, h.UsersRepository, h.CfgSrv, h.TokenSrv)
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
