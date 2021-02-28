package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AngelVlc/todos/internal/api/dtos"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

const (
	refreshTokenCookieName = "refreshToken"
)

// TokenHandler is the handler for the auth/token endpoint
func TokenHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	l, err := parseTokenBody(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	foundUser, err := h.UsersSrv.FindUserByName(l.UserName)
	if err != nil {
		return results.ErrorResult{err}
	}

	if foundUser == nil {
		return results.ErrorResult{&appErrors.BadRequestError{Msg: "The user does not exist", InternalError: nil}}
	}

	err = h.UsersSrv.CheckIfUserPasswordIsOk(foundUser, l.Password)
	if err != nil {
		return results.ErrorResult{&appErrors.BadRequestError{Msg: "Invalid password", InternalError: err}}
	}

	tokens, err := h.AuthSrv.GetTokens(foundUser)
	if err != nil {
		return results.ErrorResult{err}
	}

	addRefreshTokenCookie(w, tokens.RefreshToken)

	return results.OkResult{tokens, http.StatusOK}
}

// RefreshTokenHandler is the handler for the auth/refreshtoken endpoint
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	rt, err := getRefreshTokenCookieValue(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	rtInfo, err := h.AuthSrv.ParseRefreshToken(rt)
	if err != nil {
		return results.ErrorResult{err}
	}

	foundUser, err := h.UsersSrv.FindUserByID(rtInfo.UserID)
	if err != nil {
		return results.ErrorResult{err}
	}

	if foundUser == nil {
		return results.ErrorResult{&appErrors.BadRequestError{Msg: "The user is no longer valid", InternalError: nil}}
	}

	tokens, err := h.AuthSrv.GetTokens(foundUser)
	if err != nil {
		return results.ErrorResult{err}
	}

	addRefreshTokenCookie(w, tokens.RefreshToken)

	return results.OkResult{tokens, http.StatusOK}
}

func parseTokenBody(r *http.Request) (*dtos.TokenDto, error) {
	if r.Body == nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: nil}
	}

	decoder := json.NewDecoder(r.Body)

	var dto dtos.TokenDto
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	if len(dto.UserName) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "UserName is mandatory", InternalError: nil}
	}

	if len(dto.Password) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "Password is mandatory", InternalError: nil}
	}

	return &dto, nil
}

func addRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	rfCookie := http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/auth",
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
