package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

// UsersHandler is the handler for the users endpoints
func UsersHandler(r *http.Request, db *gorm.DB) handlerResult {
	switch r.Method {
	case http.MethodPost:
		return processUsersPOST(r, db)
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}
}

func processUsersPOST(r *http.Request, db *gorm.DB) handlerResult {
	dto, err := parseUserBody(r)
	if err != nil {
		return errorResult{err}
	}
	userSrv := wire.InitUsersService(db)
	id, err := userSrv.AddUser(&dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func parseUserBody(r *http.Request) (dtos.UserDto, error) {
	decoder := json.NewDecoder(r.Body)
	var dto dtos.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return dtos.UserDto{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return dto, nil
}
