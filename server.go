package main

import (
	"net/http"

	"github.com/AngelVlc/todos/handlers"
	"github.com/AngelVlc/todos/repositories"
	"github.com/AngelVlc/todos/services"
	"github.com/AngelVlc/todos/wire"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type server struct {
	http.Handler
	usersRepo       repositories.UsersRepository
	listsRepo       repositories.ListsRepository
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
		usersRepo:       wire.InitUsersRepository(db),
		listsRepo:       wire.InitListsRepository(db),
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
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.GetUserSingleListHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.DeleteUserListHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.UpdateUserListHandler)).Methods(http.MethodPut)
	listsSubRouter.Handle("/{listId:[0-9]+}/items", s.getHandler(handlers.AddUserListItemHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{itemId:[0-9]+}", s.getHandler(handlers.GetUserSingleListItemHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{itemId:[0-9]+}", s.getHandler(handlers.DeleteUserListItemHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{itemId:[0-9]+}", s.getHandler(handlers.UpdateUserListItemHandler)).Methods(http.MethodPut)
	listsSubRouter.Use(s.authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(handlers.AddUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Handle("", s.getHandler(handlers.GetUsersHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.GetUserHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.DeleteUserHandler)).Methods(http.MethodDelete)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.UpdateUserHandler)).Methods(http.MethodPut)
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
