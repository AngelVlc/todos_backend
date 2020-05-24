package handlers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
)

const (
	refreshTokenCookieName = "refreshToken"
)

// TokenHandler is the handler for the auth/token endpoint
func TokenHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	l, err := parseTokenBody(r)
	if err != nil {
		return errorResult{err}
	}

	foundUser, err := h.usersSrv.FindUserByName(l.UserName)
	if err != nil {
		return errorResult{err}
	}

	if foundUser == nil {
		return errorResult{&appErrors.BadRequestError{Msg: "The user does not exist", InternalError: nil}}
	}

	err = h.usersSrv.CheckIfUserPasswordIsOk(foundUser, l.Password)
	if err != nil {
		return errorResult{&appErrors.BadRequestError{Msg: "Invalid password", InternalError: err}}
	}

	tokens, err := h.authSrv.GetTokens(foundUser)
	if err != nil {
		return errorResult{err}
	}

	addRefreshTokenCookie(w, tokens["refreshToken"])

	delete(tokens, "refreshToken")

	return okResult{tokens, http.StatusOK}
}

// RefreshTokenHandler is the handler for the auth/refreshtoken endpoint
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	rt, err := getRefreshTokenCookieValue(r)
	if err != nil {
		return errorResult{err}
	}

	rtInfo, err := h.authSrv.ParseRefreshToken(rt)
	if err != nil {
		return errorResult{err}
	}

	foundUser, err := h.usersSrv.FindUserByID(rtInfo.UserID)
	if err != nil {
		return errorResult{err}
	}

	if foundUser == nil {
		return errorResult{&appErrors.BadRequestError{Msg: "The user is no longer valid", InternalError: nil}}
	}

	tokens, err := h.authSrv.GetTokens(foundUser)
	if err != nil {
		return errorResult{err}
	}

	addRefreshTokenCookie(w, tokens["refreshToken"])

	delete(tokens, "refreshToken")

	return okResult{tokens, http.StatusOK}
}

func parseTokenBody(r *http.Request) (*models.Login, error) {
	if r.Body == nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: nil}
	}

	decoder := json.NewDecoder(r.Body)

	var l models.Login
	err := decoder.Decode(&l)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	if len(l.UserName) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "UserName is mandatory", InternalError: nil}
	}

	if len(l.Password) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "Password is mandatory", InternalError: nil}
	}

	return &l, nil
}

func addRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	rfCookie := http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, &rfCookie)
}

func getRefreshTokenCookieValue(r *http.Request) (string, error) {
	c, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", &appErrors.BadRequestError{Msg: "Missing refresh token cookie", InternalError: err}
	}

	return c.Value, nil
}
