package controllers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

// TokenHandler is the handler for the auth/token endpoint
func TokenHandler(r *http.Request, db *gorm.DB) handlerResult {
	l, err := parseTokenBody(r)
	if err != nil {
		return errorResult{err}
	}

	userSrv := wire.InitUsersService(db)
	foundUser, err := userSrv.CheckIfUserPasswordIsOk(l.UserName, l.Password)
	if err != nil {
		return errorResult{err}
	}

	authSrv := wire.InitAuthService()

	tokens, err := authSrv.CreateTokens(foundUser)
	if err != nil {
		return errorResult{err}
	}

	return okResult{tokens, http.StatusOK}
}

// RefreshTokenHandler is the handler for the auth/refreshtoken endpoint
func RefreshTokenHandler(r *http.Request, db *gorm.DB) handlerResult {
	rt, err := parseRefreshTokenBody(r)
	if err != nil {
		return errorResult{err}
	}

	authSrv := wire.InitAuthService()
	rtInfo, err := authSrv.ParseRefreshToken(rt.RefreshToken)
	if err != nil {
		return errorResult{err}
	}

	userSrv := wire.InitUsersService(db)

	foundUser := userSrv.GetUserByID(rtInfo.UserID)
	if foundUser == nil {
		return errorResult{&appErrors.BadRequestError{Msg: "The user is no longer valid", InternalError: nil}}
	}

	tokens, err := authSrv.CreateTokens(foundUser)
	if err != nil {
		return errorResult{err}
	}

	return okResult{tokens, http.StatusOK}
}

func parseTokenBody(r *http.Request) (models.Login, error) {
	decoder := json.NewDecoder(r.Body)

	var l models.Login
	err := decoder.Decode(&l)
	if err != nil {
		return models.Login{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	if len(l.UserName) == 0 {
		return models.Login{}, &appErrors.BadRequestError{Msg: "UserName is mandatory", InternalError: nil}
	}

	if len(l.Password) == 0 {
		return models.Login{}, &appErrors.BadRequestError{Msg: "Password is mandatory", InternalError: nil}
	}

	return l, nil
}

func parseRefreshTokenBody(r *http.Request) (*models.RefreshToken, error) {
	decoder := json.NewDecoder(r.Body)

	var rt models.RefreshToken
	err := decoder.Decode(&rt)
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	if len(rt.RefreshToken) == 0 {
		return nil, &appErrors.BadRequestError{Msg: "RefreshToken is mandatory", InternalError: nil}
	}

	return &rt, nil
}
