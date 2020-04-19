package midlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/AngelVlc/todos/consts"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

type RequestCounterMidleware struct {
	db *gorm.DB
}

func NewRequestCounterMidleware(db *gorm.DB) RequestCounterMidleware {
	return RequestCounterMidleware{db}
}

func (m *RequestCounterMidleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := wire.InitCountersService(m.db)

		v := s.IncrementCounter("requests")

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextRequestKey, strconv.Itoa(int(v)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
