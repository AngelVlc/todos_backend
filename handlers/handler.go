package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AngelVlc/todos/consts"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	Db *gorm.DB
}

type HandlerResult interface {
	IsError() bool
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
			WriteErrorResponse(r, w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			WriteErrorResponse(r, w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			WriteErrorResponse(r, w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			WriteErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else {
			WriteErrorResponse(r, w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(okResult)
		WriteOkResponse(r, w, okRes.statusCode, okRes.content)
	}
}

// WriteErrorResponse is used when and endpoind responds with an error
func WriteErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, msg string, internalError error) {
	requestID := GetRequestIDFromContext(r)
	if internalError != nil {
		log.Printf("[%v] %v %v (%v)", requestID, statusCode, msg, internalError)
	} else {
		log.Printf("[%v] %v %v", requestID, statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

// WriteOkResponse is used when and endpoind does not respond with an error
func WriteOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v", GetRequestIDFromContext(r), statusCode)

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

func GetRequestIDFromContext(r *http.Request) string {
	requestIDRaw := r.Context().Value(consts.ReqContextRequestKey)

	requestID, _ := requestIDRaw.(string)

	return requestID
}

func parseBody(r *http.Request, dto interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return nil
}

func parseInt32UrlVar(r *http.Request, varName string) (int32, error) {
	vars := mux.Vars(r)
	value := vars[varName]
	res, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return -1, &appErrors.BadRequestError{Msg: "Invalid id in url", InternalError: err}
	}
	return int32(res), nil
}
