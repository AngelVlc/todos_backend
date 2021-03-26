package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	sharedDomain "github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/stretchr/testify/mock"
)

type RequestCounterMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

type MockedRequestCounterMiddleware struct {
	mock.Mock
}

func NewMockedRequestCounterMiddleware() *MockedRequestCounterMiddleware {
	return &MockedRequestCounterMiddleware{}
}

func (m *MockedRequestCounterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type DefaultRequestCounterMiddleware struct {
	countersRepo sharedDomain.CountersRepository
}

func NewDefaultRequestCounterMiddleware(countersRepo sharedDomain.CountersRepository) *DefaultRequestCounterMiddleware {
	return &DefaultRequestCounterMiddleware{countersRepo}
}

func (m *DefaultRequestCounterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svc := sharedApp.NewIncrementRequestsCounterService(m.countersRepo)
		v, err := svc.IncrementRequestsCounter()
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
