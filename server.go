package main

import (
	"net/http"

	"github.com/AngelVlc/todos/controllers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type server struct {
	http.Handler
	db *gorm.DB
}

func newServer(db *gorm.DB) *server {
	s := new(server)
	s.db = db

	router := mux.NewRouter()

	listsSubRouter := router.PathPrefix("/lists").Subrouter()
	listsSubRouter.Handle("", s.getHandler(controllers.GetUserLists, true, false)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(controllers.AddUserList, true, false)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.GetUserSingleList, true, false)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.DeleteUserList, true, false)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.UpdateUserList, true, false)).Methods(http.MethodPut)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(controllers.AddUserHandler, true, true)).Methods(http.MethodPost)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/token", s.getHandler(controllers.TokenHandler, false, false)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(controllers.RefreshTokenHandler, false, false)).Methods(http.MethodPost)

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
