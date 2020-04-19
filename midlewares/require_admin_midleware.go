package midlewares

import (
	"net/http"

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/controllers"
)

type RequireAdminMiddleware struct {
}

func NewRequireAdminMiddleware() RequireAdminMiddleware {
	return RequireAdminMiddleware{}
}

func (m *RequireAdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.getUserIsAdminFromContext(r) {
			controllers.WriteErrorResponse(r, w, http.StatusForbidden, "Access forbidden", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *RequireAdminMiddleware) getUserIsAdminFromContext(r *http.Request) bool {
	rawValue := r.Context().Value(consts.ReqContextUserIsAdminKey)

	isAdmin, _ := rawValue.(bool)

	return isAdmin
}
