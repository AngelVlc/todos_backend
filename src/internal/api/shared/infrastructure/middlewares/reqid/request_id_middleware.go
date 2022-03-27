package reqid

import (
	"context"
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
	"github.com/google/uuid"
)

type RequestIdMiddleware struct{}

func NewRequestIdMiddleware() *RequestIdMiddleware {
	return &RequestIdMiddleware{}
}

func (m *RequestIdMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.Header[http.CanonicalHeaderKey("X-Request-ID")]
		reqId := uuid.NewString()

		if len(values) > 0 {
			reqId = values[0]
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextRequestKey, reqId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
