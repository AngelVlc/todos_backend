package main

import (
	"net/http"

	"github.com/AngelVlc/todos/handlers"
	"github.com/AngelVlc/todos/services"
	"github.com/AngelVlc/todos/wire"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type server struct {
	http.Handler
	authSrv         services.AuthService
	listsSrv        services.ListsService
	usersSrv        services.UsersService
	authMdw         handlers.AuthMiddleware
	countersMdw     handlers.RequestCounterMiddleware
	requireAdminMdw handlers.RequireAdminMiddleware
	logMdw          handlers.LogMiddleware
}

func newServer(db *gorm.DB) *server {
	s := server{
		authSrv:         wire.InitAuthService(),
		listsSrv:        wire.InitListsService(db),
		usersSrv:        wire.InitUsersService(db),
		countersMdw:     wire.InitRequestCounterMiddleware(db),
		authMdw:         wire.InitAuthMiddleware(db),
		requireAdminMdw: wire.InitRequireAdminMiddleware(),
		logMdw:          wire.InitLogMiddleware(),
	}

	router := mux.NewRouter()
	router.Use(s.countersMdw.Middleware)

	listsSubRouter := router.PathPrefix("/lists").Subrouter()
	listsSubRouter.Handle("", s.getHandler(handlers.GetUserListsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(handlers.AddUserListHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.GetUserSingleListHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.DeleteUserListHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id}", s.getHandler(handlers.UpdateUserListHandler)).Methods(http.MethodPut)
	listsSubRouter.Use(s.authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(handlers.AddUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Use(s.authMdw.Middleware)
	usersSubRouter.Use(s.requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/token", s.getHandler(handlers.TokenHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(handlers.RefreshTokenHandler)).Methods(http.MethodPost)

	router.Use(s.logMdw.Middleware)

	s.Handler = router

	return &s
}

func (s *server) getHandler(handlerFunc handlers.HandlerFunc) handlers.Handler {
	return handlers.NewHandler(handlerFunc, s.usersSrv, s.authSrv, s.listsSrv)
}
