//+build !e2e

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}

func TestServerPublicRoutes(t *testing.T) {
	s := NewServer(nil)

	var publicRoutes = []struct {
		url            string
		method         string
		expectedStatus int
	}{
		{"/auth/login", http.MethodPost, http.StatusBadRequest},
		{"/auth/refreshtoken", http.MethodPost, http.StatusBadRequest},
	}

	for _, r := range publicRoutes {
		t.Run(fmt.Sprintf("returns %v for %v '%v'", r.expectedStatus, r.method, r.url), func(t *testing.T) {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equalf(t, r.expectedStatus, res.Result().StatusCode, "Public route %v '%v' should return %v", r.method, r.url, r.expectedStatus)
		})
	}
}

func TestServerAdminRoutes(t *testing.T) {
	s := NewServer(nil)

	var adminRoutes = []struct {
		url    string
		method string
	}{
		{"/users", http.MethodPost},
		{"/users", http.MethodGet},
		{"/users/12", http.MethodDelete},
		{"/users/12", http.MethodPut},
		{"/users/12", http.MethodGet},
	}

	for _, r := range adminRoutes {
		t.Run(fmt.Sprintf("returns 405 for %v '%v' without auth", r.url, r.method), func(t *testing.T) {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equalf(t, http.StatusUnauthorized, res.Result().StatusCode, "Admin route %v '%v' is not checking auth", r.method, r.url)
		})
	}

	for _, r := range adminRoutes {
		t.Run(fmt.Sprintf("returns 403 for %v '%v' with auth but without admin", r.url, r.method), func(t *testing.T) {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			ctx := req.Context()
			ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, false)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req.WithContext(ctx))

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, fmt.Sprintf("Admin route %v '%v' is not checking admin", r.method, r.url))
		})
	}
}

func TestServerPrivateRoutes(t *testing.T) {
	s := NewServer(nil)

	var privateRoutes = []struct {
		url    string
		method string
	}{
		{"/lists", http.MethodGet},
		{"/lists", http.MethodPost},
		{"/lists/12", http.MethodPut},
		{"/lists/12", http.MethodGet},
		{"/lists/12", http.MethodDelete},
		{"/lists/12/items", http.MethodPost},
		{"/lists/12/items/3", http.MethodGet},
		{"/lists/12/items/3", http.MethodDelete},
		{"/lists/12/items/3", http.MethodPut},
	}

	for _, r := range privateRoutes {
		t.Run(fmt.Sprintf("returns 405 for %v '%v' without auth", r.url, r.method), func(t *testing.T) {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equalf(t, http.StatusUnauthorized, res.Result().StatusCode, "Private route %v '%v' is not checking auth", r.method, r.url)
		})
	}
}

func TestServerBadRoutes(t *testing.T) {
	s := NewServer(nil)

	var badParamsRoutes = []struct {
		url    string
		method string
	}{
		{"/users/wadus", http.MethodDelete},
		{"/users/wadus", http.MethodPut},
		{"/users/wadus", http.MethodGet},
		{"/lists/wadus", http.MethodPut},
		{"/lists/wadus", http.MethodGet},
		{"/lists/wadus", http.MethodDelete},
		{"/lists/wadus/items", http.MethodPost},
		{"/lists/wadus/items/3", http.MethodGet},
		{"/lists/wadus/items/3", http.MethodDelete},
		{"/lists/wadus/items/3", http.MethodPut},
		{"/lists/3/items/wadus", http.MethodGet},
		{"/lists/3/items/wadus", http.MethodDelete},
		{"/lists/3/items/wadus", http.MethodPut},
	}

	for _, r := range badParamsRoutes {
		t.Run(fmt.Sprintf("returns 404 for %v '%v'", r.method, r.url), func(t *testing.T) {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusNotFound, res.Result().StatusCode, fmt.Sprintf("Route %v '%v' should return 404", r.method, r.url))
		})
	}
}
