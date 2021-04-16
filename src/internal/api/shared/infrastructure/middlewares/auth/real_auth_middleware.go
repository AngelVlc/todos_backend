package authmdw

import (
	"context"
	"net/http"

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
			helpers.WriteErrorResponse(r, w, http.StatusUnauthorized, "Invalid authorization token", err)
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
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", &appErrors.UnauthorizedError{Msg: "No authorization cookie", InternalError: nil}
	}

	return cookie.Value, nil
}
