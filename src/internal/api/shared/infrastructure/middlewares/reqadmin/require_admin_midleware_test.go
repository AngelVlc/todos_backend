//go:build !e2e
// +build !e2e

package reqadminmdw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/stretchr/testify/assert"
)

func TestRequireAdminMiddleware(t *testing.T) {
	md := NewRequireAdminMiddleware()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	t.Run("should return 403 if the user is not an admin", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
		assert.Equal(t, "Access forbidden\n", string(response.Body.String()))
	})

	t.Run("should call next handler when the user is an admin", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, true)
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request.WithContext(ctx))

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})
}
