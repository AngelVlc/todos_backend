package server

import (
	"net/http"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authInfra "github.com/AngelVlc/todos/internal/api/auth/infrastructure"
	listsDomain "github.com/AngelVlc/todos/internal/api/lists/domain"
	listsInfra "github.com/AngelVlc/todos/internal/api/lists/infrastructure"
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
	listsRepo     repositories.ListsRepository
	listItemsRepo repositories.ListItemsRepository
	listsSrv      services.ListsService
	listItemsSrv  services.ListItemsService
	authRepo      authDomain.AuthRepository
	listsRepoOK   listsDomain.ListsRepository
	cfgSrv        sharedApp.ConfigurationService
	passGen       authDomain.PasswordGenerator
}

func NewServer(db *gorm.DB) *server {
	s := server{
		listsRepo:     wire.InitListsRepository(db),
		listItemsRepo: wire.InitListItemsRepository(db),
		listsSrv:      wire.InitListsService(db),
		listItemsSrv:  wire.InitListItemsService(db),
		authRepo:      wire.InitAuthRepository(db),
		listsRepoOK:   wire.InitListsRepositoryOK(db),
		cfgSrv:        wire.InitConfigurationService(),
		passGen:       wire.InitPasswordGenerator(),
	}

	router := mux.NewRouter()
	countersMdw := wire.InitRequestCounterMiddleware(db)
	router.Use(countersMdw.Middleware)

	authMdw := wire.InitAuthMiddleware(db)
	requireAdminMdw := wire.InitRequireAdminMiddleware()

	listsSubRouter := router.PathPrefix("/lists").Subrouter()
	listsSubRouter.Handle("", s.getHandler(listsInfra.GetAllListsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(listsInfra.CreateListHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsInfra.GetListHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsInfra.DeleteListHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsInfra.UpdateListHandler)).Methods(http.MethodPut)
	listsSubRouter.Handle("/{listId:[0-9]+}/items", s.getHandler(listsInfra.GetAllListItemsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{listId:[0-9]+}/items", s.getHandler(listsInfra.CreateListItemHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsInfra.GetListItemHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsInfra.DeleteListItemHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsInfra.UpdateListItemHandler)).Methods(http.MethodPut)
	listsSubRouter.Use(authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(authInfra.CreateUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Handle("", s.getHandler(authInfra.GetAllUsersHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authInfra.GetUserHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authInfra.DeleteUserHandler)).Methods(http.MethodDelete)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authInfra.UpdateUserHandler)).Methods(http.MethodPut)
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
	return handler.NewHandler(handlerFunc, s.listsSrv, s.listItemsSrv, s.authRepo, s.listsRepoOK, s.cfgSrv, s.passGen)
}
