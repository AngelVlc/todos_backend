//go:build !e2e
// +build !e2e

package authmdw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AngelVlc/todos_backend/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/consts"
	"github.com/golang-jwt/jwt"

	"github.com/stretchr/testify/assert"
)

func TestRealAuthMiddleware(t *testing.T) {
	mockedTokenSrv := domain.NewMockedTokenService()
	md := NewRealAuthMiddleware(mockedTokenSrv)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := ctx.Value(consts.ReqContextUserIDKey).(int32)
		userName, _ := ctx.Value(consts.ReqContextUserNameKey).(string)
		isAdmin, _ := ctx.Value(consts.ReqContextUserIsAdminKey).(bool)

		assert.Equal(t, int32(1), userID)
		assert.Equal(t, "user", userName)
		assert.True(t, isAdmin)
	})

	getTokenCookie := func(rt string) *http.Cookie {
		return &http.Cookie{Name: "token", Value: rt}
	}

	t.Run("Should return an error if there isn't token cookie", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "No authorization cookie\n", string(response.Body.String()))
	})

	t.Run("Should return an error if the token is not valid", func(t *testing.T) {
		mockedTokenSrv.On("ParseToken", "badToken").Return(nil, fmt.Errorf("some error")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.AddCookie(getTokenCookie("badToken"))
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid authorization token\n", string(response.Body.String()))

		mockedTokenSrv.AssertExpectations(t)
	})

	t.Run("Should add the token info to the request context if the token is valid", func(t *testing.T) {
		token := jwt.Token{Valid: true}
		mockedTokenSrv.On("ParseToken", "validToken").Return(&token, nil).Once()
		mockedTokenSrv.On("GetTokenInfo", &token).Return(&domain.TokenClaimsInfo{UserID: 1, UserName: "user", IsAdmin: true}).Once()

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.AddCookie(getTokenCookie("validToken"))
		response := httptest.NewRecorder()
		handlerToTest := md.Middleware(nextHandler)

		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		mockedTokenSrv.AssertExpectations(t)
	})
}
