package middleware

import (
	"net/http"
)

type AuthMiddleware interface {
	Middleware(next http.Handler) http.Handler
}