package reqid

import (
	"context"
	"net/http"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/google/uuid"
)

type RequestIdMiddleware struct{}

func NewRequestIdMiddleware() *RequestIdMiddleware {
	return &RequestIdMiddleware{}
}

func (m *RequestIdMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextRequestKey, uuid.NewString())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
