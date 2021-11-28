package recover

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/honeybadger-io/honeybadger-go"
)

type RecoverMiddleware struct{}

func NewRecoverMiddleware() *RecoverMiddleware {
	return &RecoverMiddleware{}
}

func (m *RecoverMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				honeybadger.Notify(err)
				log.Println(string(debug.Stack()))

				helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, "Internal error", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
