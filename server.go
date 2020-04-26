package main

import (
	"net/http"

	"github.com/AngelVlc/todos/controllers"
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
	listsSubRouter.Handle("", s.getHandler(controllers.GetUserLists)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(controllers.AddUserList)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.GetUserSingleList)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.DeleteUserList)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id}", s.getHandler(controllers.UpdateUserList)).Methods(http.MethodPut)
	listsSubRouter.Use(authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(controllers.AddUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Use(authMdw.Middleware)
	usersSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/token", s.getHandler(controllers.TokenHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(controllers.RefreshTokenHandler)).Methods(http.MethodPost)

	logMdw := midlewares.NewLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc) controllers.Handler {
	return controllers.Handler{
		HandlerFunc: handlerFunc,
		Db:          s.db,
	}
}
