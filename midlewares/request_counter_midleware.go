package midlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/controllers"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

type RequestCounterMiddleware struct {
	db *gorm.DB
}

func NewRequestCounterMiddleware(db *gorm.DB) RequestCounterMiddleware {
	return RequestCounterMiddleware{db}
}

func (m *RequestCounterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := wire.InitCountersService(m.db)

		v, err := s.IncrementCounter("requests")
		if err != nil {
			log.Printf("[] %v %q", r.Method, r.URL)
			controllers.WriteErrorResponse(r, w, http.StatusInternalServerError, "Error incrementing requests counter", err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextRequestKey, strconv.Itoa(int(v)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
