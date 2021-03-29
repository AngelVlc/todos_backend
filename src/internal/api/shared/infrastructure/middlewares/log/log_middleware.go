package logmdw

import (
	"log"
	"net/http"

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
		userName := m.getUserNameFromContext(r)
		if len(userName) > 0 {
			log.Printf("[%v] %v %q", requestID, r.Method, r.URL)
		} else {
			log.Printf("[%v] %v %v %q", requestID, userName, r.Method, r.URL)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *LogMiddleware) getUserNameFromContext(r *http.Request) string {
	rawValue := r.Context().Value(consts.ReqContextUserNameKey)

	name, _ := rawValue.(string)

	return name
}
