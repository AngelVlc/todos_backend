package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	listsDomain "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/gorilla/mux"
	"github.com/honeybadger-io/honeybadger-go"
	"gorm.io/gorm"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	AuthRepository       authDomain.AuthRepository
	UsersRepository      authDomain.UsersRepository
	ListsRepository      listsDomain.ListsRepository
	CategoriesRepository listsDomain.CategoriesRepository
	CfgSrv               sharedApp.ConfigurationService
	TokenSrv             authDomain.TokenService
	PassGen              passgen.PasswordGenerator
	EventBus             events.EventBus
	RequestInput         interface{}
	SearchClient         search.SearchIndexClient
}

type HandlerResult interface {
	IsError() bool
}

func NewHandler(f HandlerFunc,
	authRepo authDomain.AuthRepository,
	usersRepo authDomain.UsersRepository,
	listsRepo listsDomain.ListsRepository,
	categoriesRepo listsDomain.CategoriesRepository,
	cfgSrv sharedApp.ConfigurationService,
	tokenSrv authDomain.TokenService,
	passGen passgen.PasswordGenerator,
	eventBus events.EventBus,
	requestInput interface{},
	searchClient search.SearchIndexClient) Handler {

	return Handler{
		HandlerFunc:          f,
		AuthRepository:       authRepo,
		UsersRepository:      usersRepo,
		ListsRepository:      listsRepo,
		CategoriesRepository: categoriesRepo,
		CfgSrv:               cfgSrv,
		PassGen:              passGen,
		TokenSrv:             tokenSrv,
		EventBus:             eventBus,
		RequestInput:         requestInput,
		SearchClient:         searchClient,
	}
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request, Handler) HandlerResult

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.RequestInput != nil {
		err := h.parseBody(r, h.RequestInput)
		if err != nil {
			if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
				helpers.WriteErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
			} else {
				helpers.WriteErrorResponse(r, w, http.StatusBadRequest, "Invalid body", err)
			}
			return
		}
	}

	res := h.HandlerFunc(w, r, h)

	if res.IsError() {
		errorRes, _ := res.(results.ErrorResult)
		err := errorRes.Err
		if unexErr, ok := err.(*appErrors.UnexpectedError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			helpers.WriteErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.WriteErrorResponse(r, w, http.StatusNotFound, "Not found", err)
		} else {
			honeybadger.Notify(err)
			helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(results.OkResult)
		helpers.WriteOkResponse(r, w, okRes.StatusCode, okRes.Content)
	}
}

func (h Handler) parseBody(r *http.Request, result interface{}) error {
	if r.Body == nil {
		return &appErrors.BadRequestError{Msg: "Invalid body"}
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			return badRequestErr
		} else {
			return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
		}
	}

	return nil
}

func (h Handler) ParseInt32UrlVar(r *http.Request, varName string) int32 {
	vars := mux.Vars(r)
	value := vars[varName]
	res, _ := strconv.ParseInt(value, 10, 32)

	return int32(res)
}

func (h Handler) GetUserIDFromContext(r *http.Request) int32 {
	userIDRaw := r.Context().Value(consts.ReqContextUserIDKey)

	userID, _ := userIDRaw.(int32)

	return userID
}
