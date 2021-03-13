//+build !e2e

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"

	"github.com/stretchr/testify/assert"
)

func TestRealAuthMiddleware(t *testing.T) {

	mockedCfgSrv := sharedApp.NewMockedConfigurationService()
	md := NewRealAuthMiddleware(mockedCfgSrv)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := ctx.Value(consts.ReqContextUserIDKey).(int32)
		userName, _ := ctx.Value(consts.ReqContextUserNameKey).(string)
		isAdmin, _ := ctx.Value(consts.ReqContextUserIsAdminKey).(bool)

		assert.Equal(t, int32(1), userID)
		assert.Equal(t, "user", userName)
		assert.True(t, isAdmin)
	})

	t.Run("Should return an error if there isn't authorization header", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "No authorization header\n", string(response.Body.String()))
	})

	t.Run("Should return an error if the authorization header is not a bearer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "bad_header")
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid authorization header\n", string(response.Body.String()))
	})

	t.Run("Should return an error if the token is not a jwt token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer badToken")
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		mockedCfgSrv.On("GetJwtSecret").Return("secret").Once()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Error parsing the authorization token\n", string(response.Body.String()))
		mockedCfgSrv.AssertExpectations(t)
	})

	t.Run("Should add the token info to the request context if the token is valid", func(t *testing.T) {
		mockedCfgSrv.On("TokenExpirationInSeconds").Return(5 * time.Minute).Once()
		mockedCfgSrv.On("GetJwtSecret").Return("secret").Times(2)
		authUSer := domain.User{ID: int32(1), Name: "user", IsAdmin: true}
		token, _ := domain.NewTokenService(mockedCfgSrv).GenerateToken(&authUSer)
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})
}
