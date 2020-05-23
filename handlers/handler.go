package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AngelVlc/todos/consts"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/services"
	"github.com/gorilla/mux"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	usersSrv services.UsersService
	authSrv  services.AuthService
	listsSrv services.ListsService
}

type HandlerResult interface {
	IsError() bool
}

func NewHandler(f HandlerFunc, u services.UsersService, a services.AuthService, l services.ListsService) Handler {
	return Handler{
		HandlerFunc: f,
		usersSrv:    u,
		authSrv:     a,
		listsSrv:    l,
	}
}

type errorResult struct {
	err error
}

func (e errorResult) IsError() bool {
	return true
}

type okResult struct {
	content    interface{}
	statusCode int
}

func (r okResult) IsError() bool {
	return false
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(*http.Request, Handler) HandlerResult

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.HandlerFunc(r, h)

	if res.IsError() {
		errorRes, _ := res.(errorResult)
		err := errorRes.err
		if unexErr, ok := err.(*appErrors.UnexpectedError); ok {
			writeErrorResponse(r, w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			writeErrorResponse(r, w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			writeErrorResponse(r, w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			writeErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else {
			writeErrorResponse(r, w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(okResult)
		writeOkResponse(r, w, okRes.statusCode, okRes.content)
	}
}

// writeErrorResponse is used when and endpoind responds with an error
func writeErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, msg string, internalError error) {
	requestID := getRequestIDFromContext(r)
	if internalError != nil {
		log.Printf("[%v] %v %v (%v)", requestID, statusCode, msg, internalError)
	} else {
		log.Printf("[%v] %v %v", requestID, statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

// writeOkResponse is used when and endpoind does not respond with an error
func writeOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v", getRequestIDFromContext(r), statusCode)

	const jsonContentType = "application/json"

	if content != nil {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(content)
	} else {
		w.WriteHeader(statusCode)
	}
}

func getUserIDFromContext(r *http.Request) int32 {
	userIDRaw := r.Context().Value(consts.ReqContextUserIDKey)

	userID, _ := userIDRaw.(int32)

	return userID
}

func getRequestIDFromContext(r *http.Request) string {
	requestIDRaw := r.Context().Value(consts.ReqContextRequestKey)

	requestID, _ := requestIDRaw.(string)

	return requestID
}

func parseBody(r *http.Request, dto interface{}) error {
	if r.Body == nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: nil}
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return nil
}

func parseInt32UrlVar(r *http.Request, varName string) int32 {
	vars := mux.Vars(r)
	value := vars[varName]
	res, _ := strconv.ParseInt(value, 10, 32)
	return int32(res)
}
