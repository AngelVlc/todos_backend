package reqcountermdw

import (
	"context"
	"log"
	"net/http"
	"strconv"

	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
)

type RequestCounterMiddleware struct {
	incrementRequestsCounterSvc *sharedApp.IncrementRequestsCounterService
}

func NewRequestCounterMiddleware(incReqCounterSrv *sharedApp.IncrementRequestsCounterService) *RequestCounterMiddleware {
	return &RequestCounterMiddleware{incReqCounterSrv}
}

func (m *RequestCounterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := m.incrementRequestsCounterSvc.IncrementRequestsCounter()
		if err != nil {
			log.Printf("[] %v %q", r.Method, r.URL)
			helpers.WriteErrorResponse(r, w, http.StatusInternalServerError, "Error incrementing requests counter", err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.ReqContextRequestKey, strconv.Itoa(int(v)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
