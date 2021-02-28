package handler

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/services"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	UsersSrv     services.UsersService
	AuthSrv      services.AuthService
	ListsSrv     services.ListsService
	ListItemsSrv services.ListItemsService
}

type HandlerResult interface {
	IsError() bool
}

func NewHandler(f HandlerFunc, usersSvc services.UsersService, authSvc services.AuthService, listsSvc services.ListsService, listItemsSvc services.ListItemsService) Handler {
	return Handler{
		HandlerFunc:  f,
		UsersSrv:     usersSvc,
		AuthSrv:      authSvc,
		ListsSrv:     listsSvc,
		ListItemsSrv: listItemsSvc,
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
