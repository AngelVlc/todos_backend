package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/wire"

	// "github.com/AngelVlc/todos/services"
	"github.com/jinzhu/gorm"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	Db           *gorm.DB
	RequireAuth  bool
	RequireAdmin bool
}

type handlerResult interface {
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
type HandlerFunc func(*http.Request, *gorm.DB) handlerResult

type contextKey string

const reqContextUserKey contextKey = "userID"
const reqContextRequestKey contextKey = "requestID"

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jwtInfo *models.JwtClaimsInfo

	r = h.addRequestIDToContext(r)

	if h.RequireAuth {
		token, err := h.getAuthToken(r)
		if err != nil {
			h.writeErrorResponse(r, w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		authSrv := wire.InitAuthService()
		jwtInfo, err = authSrv.ParseToken(token)
		if err != nil {
			h.writeErrorResponse(r, w, http.StatusUnauthorized, "Invalid auth token", err)
			return
		}

		if h.RequireAdmin && !jwtInfo.IsAdmin {
			h.writeErrorResponse(r, w, http.StatusForbidden, "Access forbidden", err)
			return
		}
	}

	requestID := h.getRequestIDFromContext(r)
	if jwtInfo == nil {
		log.Printf("[%v] %v %q", requestID, r.Method, r.URL)
	} else {
		log.Printf("[%v] %v %v %q", requestID, jwtInfo.UserName, r.Method, r.URL)

		r = h.addUserIDToContext(jwtInfo.UserID, r)
	}

	res := h.HandlerFunc(r, h.Db)

	if res.IsError() {
		errorRes, _ := res.(errorResult)
		err := errorRes.err
		if unexErr, ok := err.(*appErrors.UnexpectedError); ok {
			h.writeErrorResponse(r, w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			h.writeErrorResponse(r, w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			h.writeErrorResponse(r, w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			h.writeErrorResponse(r, w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else {
			h.writeErrorResponse(r, w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(okResult)
		h.writeOkResponse(r, w, okRes.statusCode, okRes.content)
	}
}

// writeErrorResponse is used when and endpoind responds with an error
func (h Handler) writeErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, msg string, internalError error) {
	requestID := h.getRequestIDFromContext(r)
	if internalError != nil {
		log.Printf("[%v] %v %v (%v)", requestID, statusCode, msg, internalError)
	} else {
		log.Printf("[%v] %v %v", requestID, statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

// writeOkResponse is used when and endpoind does not respond with an error
func (h Handler) writeOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v", h.getRequestIDFromContext(r), statusCode)

	const jsonContentType = "application/json"

	if content != nil {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(content)
	} else {
		w.WriteHeader(statusCode)
	}
}

func (h Handler) getAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if len(authHeader) == 0 {
		return "", &appErrors.UnauthorizedError{Msg: "No authorization header", InternalError: nil}
	}

	authHeaderParts := strings.Split(authHeader, "Bearer ")

	if len(authHeaderParts) != 2 {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid authorization header", InternalError: nil}
	}

	return authHeaderParts[1], nil
}

func (h Handler) addUserIDToContext(userID int32, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, reqContextUserKey, userID)

	return r.WithContext(ctx)
}

func (h Handler) addRequestIDToContext(r *http.Request) *http.Request {
	s := wire.InitCountersService(h.Db)

	v := s.IncrementCounter("requests")

	ctx := r.Context()
	ctx = context.WithValue(ctx, reqContextRequestKey, strconv.Itoa(int(v)))

	return r.WithContext(ctx)
}

func getUserIDFromContext(r *http.Request) int32 {
	userIDRaw := r.Context().Value(reqContextUserKey)

	userID, _ := userIDRaw.(int32)

	return userID
}

func (h Handler) getRequestIDFromContext(r *http.Request) string {
	requestIDRaw := r.Context().Value(reqContextRequestKey)

	requestID, _ := requestIDRaw.(string)

	return requestID
}
