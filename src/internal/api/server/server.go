package server

import (
	"net/http"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authInfra "github.com/AngelVlc/todos/internal/api/auth/infrastructure"
	"github.com/AngelVlc/todos/internal/api/handlers"
	"github.com/AngelVlc/todos/internal/api/repositories"
	"github.com/AngelVlc/todos/internal/api/services"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/wire"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type server struct {
	http.Handler
	usersRepo     repositories.UsersRepository
	listsRepo     repositories.ListsRepository
	listItemsRepo repositories.ListItemsRepository
	listsSrv      services.ListsService
	listItemsSrv  services.ListItemsService
	usersSrv      services.UsersService
	authRepo      authDomain.AuthRepository
	cfgSrv        sharedApp.ConfigurationService
}

func NewServer(db *gorm.DB) *server {
	s := server{
		usersRepo:     wire.InitUsersRepository(db),
		listsRepo:     wire.InitListsRepository(db),
		listItemsRepo: wire.InitListItemsRepository(db),
		listsSrv:      wire.InitListsService(db),
		listItemsSrv:  wire.InitListItemsService(db),
		usersSrv:      wire.InitUsersService(db),
		authRepo:      wire.InitAuthRepository(db),
		cfgSrv:        wire.InitConfigurationService(),
	}

	router := mux.NewRouter()
	countersMdw := wire.InitRequestCounterMiddleware(db)
	router.Use(countersMdw.Middleware)

	authMdw := wire.InitAuthMiddleware(db)
	requireAdminMdw := wire.InitRequireAdminMiddleware()

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
	listsSubRouter.Use(authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(handlers.AddUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Handle("", s.getHandler(handlers.GetUsersHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.GetUserHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.DeleteUserHandler)).Methods(http.MethodDelete)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(handlers.UpdateUserHandler)).Methods(http.MethodPut)
	usersSubRouter.Use(authMdw.Middleware)
	usersSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/login", s.getHandler(authInfra.LoginHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(authInfra.RefreshTokenHandler)).Methods(http.MethodPost)

	logMdw := wire.InitLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	return &s
}

func (s *server) getHandler(handlerFunc handler.HandlerFunc) handler.Handler {
	return handler.NewHandler(handlerFunc, s.usersSrv, s.listsSrv, s.listItemsSrv, s.authRepo, s.cfgSrv)
}
