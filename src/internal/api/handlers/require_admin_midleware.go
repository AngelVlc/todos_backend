package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/stretchr/testify/mock"
)

type RequireAdminMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

type MockedRequireAdminMiddleware struct {
	mock.Mock
}

type DefaultRequireAdminMiddleware struct {
}

func NewDefaultRequireAdminMiddleware() *DefaultRequireAdminMiddleware {
	return &DefaultRequireAdminMiddleware{}
}

func (m *DefaultRequireAdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.getUserIsAdminFromContext(r) {
			helpers.WriteErrorResponse(r, w, http.StatusForbidden, "Access forbidden", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *DefaultRequireAdminMiddleware) getUserIsAdminFromContext(r *http.Request) bool {
	rawValue := r.Context().Value(consts.ReqContextUserIsAdminKey)

	isAdmin, _ := rawValue.(bool)

	return isAdmin
}
