package handler

import (
	"encoding/json"
	"net/http"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/services"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	ListsSrv       services.ListsService
	ListItemsSrv   services.ListItemsService
	AuthRepository authDomain.AuthRepository
	CfgSrv         sharedApp.ConfigurationService
	PassGen        authDomain.PasswordGenerator
}

type HandlerResult interface {
	IsError() bool
}

func NewHandler(f HandlerFunc,
	listsSvc services.ListsService, listItemsSvc services.ListItemsService,
	authRepo authDomain.AuthRepository, cfgSrv sharedApp.ConfigurationService,
	passGen authDomain.PasswordGenerator) Handler {
	return Handler{
		HandlerFunc:    f,
		ListsSrv:       listsSvc,
		ListItemsSrv:   listItemsSvc,
		AuthRepository: authRepo,
		CfgSrv:         cfgSrv,
		PassGen:        passGen,
	}
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request, Handler) HandlerResult

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.HandlerFunc(w, r, h)

	if res.IsError() {
		errorRes, _ := res.(results.ErrorResult)
		err := errorRes.Err
		if unexErr, ok := err.(*appErrors.UnexpectedError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else {
			helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(results.OkResult)
		helpers.WriteOkResponse(r, w, okRes.StatusCode, okRes.Content)
	}
}

func (h Handler) ParseBody(r *http.Request, result interface{}) error {
	if r.Body == nil {
		return &appErrors.BadRequestError{Msg: "Invalid body"}
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return nil
}
