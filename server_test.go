package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/consts"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	server := newServer(nil)

	t.Run("handles /users without auth", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})

	t.Run("handles /users with auth but without admin", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		request.Header.Set("Authorization", "bearer")
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, false)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request.WithContext(ctx))

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
	})

	t.Run("handles /users with auth and admin", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		request.Header.Set("Authorization", "bearer")
		ctx := request.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, true)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request.WithContext(ctx))

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
	})
}
