//+build !e2e

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos/internal/api/consts"
	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/AngelVlc/todos/internal/api/services"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {

	mockedAuthSvc := services.NewMockedAuthService()
	md := NewDefaultAuthMiddleware(mockedAuthSvc)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := ctx.Value(consts.ReqContextUserIDKey).(int32)
		userName, _ := ctx.Value(consts.ReqContextUserNameKey).(string)
		isAdmin, _ := ctx.Value(consts.ReqContextUserIsAdminKey).(bool)

		assert.Equal(t, int32(11), userID)
		assert.Equal(t, "user", userName)
		assert.True(t, isAdmin)
	})

	t.Run("Should return an error if there isn't authorization header", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "No authorization header\n", string(response.Body.String()))
	})

	t.Run("Should return an error if the authorization header is not a bearer", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "bad_header")
		response := httptest.NewRecorder()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid authorization header\n", string(response.Body.String()))
	})

	t.Run("Should return an error if the token is not valid", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer badToken")
		response := httptest.NewRecorder()

		mockedAuthSvc.On("ParseToken", "badToken").Return(nil, fmt.Errorf("some error")).Once()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid authorization token\n", string(response.Body.String()))
	})

	t.Run("Should add the token info to the request context if the token is valid", func(t *testing.T) {
		handlerToTest := md.Middleware(nextHandler)

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer token")
		response := httptest.NewRecorder()

		jwtInfo := models.JwtClaimsInfo{
			UserID:   11,
			UserName: "user",
			IsAdmin:  true,
		}

		mockedAuthSvc.On("ParseToken", "token").Return(&jwtInfo, nil).Once()

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})
}
