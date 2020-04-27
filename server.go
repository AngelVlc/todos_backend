package main

import (
	"net/http"

	"github.com/AngelVlc/todos/handlers"
	"github.com/AngelVlc/todos/midlewares"
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

	countersMdw := midlewares.NewRequestCounterMiddleware(s.db)
	router.Use(countersMdw.Middleware)

	authMdw := midlewares.NewAuthMiddleware()
	requireAdminMdw := midlewares.NewRequireAdminMiddleware()

	listsSubRouter := router.PathPrefix("/lists").Subrouter()
	listsSubRouter.Handle("", s.getHandler(handlers.GetUserListsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(handlers.AddUserListHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.GetUserSingleListHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.DeleteUserListHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.UpdateUserListHandler)).Methods(http.MethodPut)
	listsSubRouter.Use(authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(handlers.AddUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Use(authMdw.Middleware)
	usersSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/token", s.getHandler(handlers.TokenHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(handlers.RefreshTokenHandler)).Methods(http.MethodPost)

	logMdw := midlewares.NewLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc handlers.HandlerFunc) handlers.Handler {
	return handlers.NewHandler(handlerFunc, s.db)
}
