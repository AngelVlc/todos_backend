package main

import (
	"net/http"

	"github.com/AngelVlc/todos/controllers"
	"github.com/jinzhu/gorm"
)

type server struct {
	http.Handler
	db *gorm.DB
}

func newServer(db *gorm.DB) *server {
	s := new(server)
	s.db = db

	router := http.NewServeMux()

	router.Handle("/lists", s.getHandler(controllers.ListsHandler, true, false))
	router.Handle("/lists/", s.getHandler(controllers.ListsHandler, true, false))
	router.Handle("/users", s.getHandler(controllers.UsersHandler, true, true))
	router.Handle("/users/", s.getHandler(controllers.UsersHandler, true, true))
	router.Handle("/auth/token", s.getHandler(controllers.TokenHandler, false, false))
	router.Handle("/auth/refreshtoken", s.getHandler(controllers.RefreshTokenHandler, false, false))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc, requireAuth bool, requireAdmin bool) controllers.Handler {
	return controllers.Handler{
		HandlerFunc:  handlerFunc,
		Db:           s.db,
		RequireAuth:  requireAuth,
		RequireAdmin: requireAdmin,
	}
}
