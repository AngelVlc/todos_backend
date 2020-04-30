package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/AngelVlc/todos/consts"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/services"
	"github.com/stretchr/testify/mock"
)

type AuthMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

type MockedAuthMiddleware struct {
	mock.Mock
}

func NewMockedAuthMiddleware() *MockedAuthMiddleware {
	return &MockedAuthMiddleware{}
}

func (m *MockedAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			writeErrorResponse(r, w, http.StatusUnauthorized, "No authorization header", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type DefaultAuthMiddleware struct {
	auth services.AuthService
}

func NewDefaultAuthMiddleware(auth services.AuthService) *DefaultAuthMiddleware {
	return &DefaultAuthMiddleware{auth}
}

func (m *DefaultAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.getAuthToken(r)
		if err != nil {
			writeErrorResponse(r, w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		jwtInfo, err := m.auth.ParseToken(token)
		if err != nil {
			writeErrorResponse(r, w, http.StatusUnauthorized, "Invalid authorization token", err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, jwtInfo.UserID)
		ctx = context.WithValue(ctx, consts.ReqContextUserNameKey, jwtInfo.UserName)
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, jwtInfo.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *DefaultAuthMiddleware) getAuthToken(r *http.Request) (string, error) {
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
