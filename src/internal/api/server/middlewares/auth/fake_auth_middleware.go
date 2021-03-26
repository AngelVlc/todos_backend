package middleware

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/stretchr/testify/mock"
)

type FakeAuthMiddleware struct {
	mock.Mock
}

func NewFakeAuthMiddleware() *FakeAuthMiddleware {
	return &FakeAuthMiddleware{}
}

func (m *FakeAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, "No authorization header", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}