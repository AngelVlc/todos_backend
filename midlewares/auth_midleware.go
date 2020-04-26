package midlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AngelVlc/todos/consts"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/handlers"
	"github.com/AngelVlc/todos/services"
	"github.com/AngelVlc/todos/wire"
)

type AuthMiddleware struct {
	auth services.AuthService
}

func NewAuthMiddleware() AuthMiddleware {
	return AuthMiddleware{wire.InitAuthService()}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.getAuthToken(r)
		if err != nil {
			handlers.WriteErrorResponse(r, w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		jwtInfo, err := m.auth.ParseToken(token)
		if err != nil {
			handlers.WriteErrorResponse(r, w, http.StatusUnauthorized, "Invalid authorization token", err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, jwtInfo.UserID)
		ctx = context.WithValue(ctx, consts.ReqContextUserNameKey, jwtInfo.UserName)
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, jwtInfo.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) getAuthToken(r *http.Request) (string, error) {
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
