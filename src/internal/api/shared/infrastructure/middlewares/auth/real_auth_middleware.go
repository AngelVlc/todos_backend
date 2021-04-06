package authmdw

import (
	"context"
	"net/http"
	"strings"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
)

type RealAuthMiddleware struct {
	tokenSrv domain.TokenService
}

func NewRealAuthMiddleware(tokenSrv domain.TokenService) *RealAuthMiddleware {
	return &RealAuthMiddleware{tokenSrv}
}

func (m *RealAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.getAuthToken(r)
		if err != nil {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		parsedToken, err := m.tokenSrv.ParseToken(token)
		if err != nil {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, "Error parsing the authorization token", err)
			return
		}

		if !parsedToken.Valid {
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, "Invalid authorization token", nil)
			return
		}

		tokenInfo := m.tokenSrv.GetTokenInfo(parsedToken)

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextUserIDKey, tokenInfo.UserID)
		ctx = context.WithValue(ctx, consts.ReqContextUserNameKey, tokenInfo.UserName)
		ctx = context.WithValue(ctx, consts.ReqContextUserIsAdminKey, tokenInfo.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *RealAuthMiddleware) getAuthToken(r *http.Request) (string, error) {
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
