package fakemdw

import (
	"net/http"
)

type FakeMiddleware struct{}

func NewFakeMiddleware() *FakeMiddleware {
	return &FakeMiddleware{}
}

func (m *FakeMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
