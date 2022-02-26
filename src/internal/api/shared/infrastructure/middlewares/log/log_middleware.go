package logmdw

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
)

type LogMiddleware struct{}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (m *LogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := helpers.GetRequestIDFromContext(r)

		log.Printf("[%v] %v %q", requestID, r.Method, r.URL)

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextStartTime, time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
