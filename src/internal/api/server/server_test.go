//+build !e2e

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AngelVlc/todos/internal/api/consts"
	appErrors "github.com/AngelVlc/todos/internal/api/errors"
	"github.com/AngelVlc/todos/internal/api/services"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}

type route struct {
	url    string
	method string
}

var publicRoutes = []route{
	{"/auth/token", http.MethodPost},
	{"/auth/refreshtoken", http.MethodPost},
}

var privateRoutes = []route{
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

var adminRoutes = []route{
	{"/users", http.MethodPost},
	{"/users", http.MethodGet},
	{"/users/12", http.MethodDelete},
	{"/users/12", http.MethodPut},
	{"/users/12", http.MethodGet},
}

var badParamsRoutes = []route{
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

func TestServer(t *testing.T) {
	s := NewServer(nil)

	t.Run("handles public routes", func(t *testing.T) {
		for _, r := range publicRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, fmt.Sprintf("Public route %v '%v' not working", r.method, r.url))
		}
	})

	t.Run("handles admin routes without auth", func(t *testing.T) {
		for _, r := range adminRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, fmt.Sprintf("Admin route %v '%v' is not checking auth", r.method, r.url))
		}
	})

	t.Run("handles admin routes with auth but without admin", func(t *testing.T) {
		for _, r := range adminRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			ctx := req.Context()
			ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, false)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req.WithContext(ctx))

			assert.Equal(t, http.StatusForbidden, res.Result().StatusCode, fmt.Sprintf("Admin route %v '%v' is not checking admin", r.method, r.url))
		}
	})

	t.Run("handles admin routes with auth and admin", func(t *testing.T) {
		err := appErrors.BadRequestError{Msg: "Some error"}
		mockedUsersSrv, _ := s.usersSrv.(*services.MockedUsersService)
		mockedUsersSrv.On("GetUsers").Return(nil, &err).Once()
		mockedUsersSrv.On("FindUserByID", int32(12)).Return(nil, &err).Once()
		mockedListsSrv, _ := s.listsSrv.(*services.MockedListsService)
		mockedListsSrv.On("GetUserLists", int32(12)).Return(nil, &err).Once()

		for _, r := range adminRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			ctx := req.Context()
			ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, true)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req.WithContext(ctx))

			assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, fmt.Sprintf("Admin route %v '%v' is not working", r.method, r.url))
		}

		mockedUsersSrv.AssertExpectations(t)
		mockedListsSrv.AssertExpectations(t)
	})

	t.Run("handles private routes without auth", func(t *testing.T) {
		for _, r := range privateRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, fmt.Sprintf("Private route %v '%v' is not checking auth", r.method, r.url))
		}
	})

	t.Run("handles private routes with auth", func(t *testing.T) {
		err := appErrors.BadRequestError{Msg: "Some error"}
		mockedListsSrv, _ := s.listsSrv.(*services.MockedListsService)
		mockedListsSrv.On("GetUserLists", int32(0)).Return(nil, &err).Once()
		mockedListsSrv.On("GetUserList", int32(12), int32(0)).Return(nil, &err).Once()
		mockedListsSrv.On("RemoveUserList", int32(12), int32(0)).Return(&err).Once()
		mockedListItemsSrv, _ := s.listItemsSrv.(*services.MockedListItemsService)
		mockedListItemsSrv.On("GetListItem", int32(3), int32(12), int32(0)).Return(nil, &err).Once()
		mockedListItemsSrv.On("RemoveListItem", int32(3), int32(12), int32(0)).Return(&err).Once()

		for _, r := range privateRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, fmt.Sprintf("Private route %v '%v' is not working", r.method, r.url))
		}

		mockedListsSrv.AssertExpectations(t)
	})

	t.Run("handles routes with bad url params", func(t *testing.T) {
		for _, r := range badParamsRoutes {
			req, _ := http.NewRequest(r.method, r.url, nil)
			req.Header.Set("Authorization", "bearer")
			res := httptest.NewRecorder()

			s.ServeHTTP(res, req)

			assert.Equal(t, http.StatusNotFound, res.Result().StatusCode, fmt.Sprintf("Route %v '%v' should return a 404 status", r.method, r.url))
		}
	})
}
