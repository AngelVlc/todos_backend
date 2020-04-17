package controllers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
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

func parseUserBody(r *http.Request) (models.UserDto, error) {
	decoder := json.NewDecoder(r.Body)
	var dto models.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.UserDto{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return dto, nil
}
