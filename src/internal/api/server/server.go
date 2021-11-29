package server

import (
	"net/http"
	"net/http/pprof"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/auth/domain/passgen"
	authInfra "github.com/AngelVlc/todos/internal/api/auth/infrastructure"
	listsDomain "github.com/AngelVlc/todos/internal/api/lists/domain"
	listsInfra "github.com/AngelVlc/todos/internal/api/lists/infrastructure"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/middlewares/recover"
	"github.com/AngelVlc/todos/internal/api/wire"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type server struct {
	http.Handler
	authRepo    authDomain.AuthRepository
	listsRepo   listsDomain.ListsRepository
	cfgSrv      sharedApp.ConfigurationService
	tokenSrv    authDomain.TokenService
	passGen     passgen.PasswordGenerator
	eventBus    events.EventBus
	subscribers []events.Subscriber
	newRelicApp *newrelic.Application
}

func NewServer(db *gorm.DB, eb events.EventBus, newRelicApp *newrelic.Application) *server {

	s := server{
		authRepo:    wire.InitAuthRepository(db),
		listsRepo:   wire.InitListsRepository(db),
		cfgSrv:      wire.InitConfigurationService(),
		tokenSrv:    wire.InitTokenService(),
		passGen:     wire.InitPasswordGenerator(),
		eventBus:    eb,
		subscribers: []events.Subscriber{},
		newRelicApp: newRelicApp,
	}

	router := mux.NewRouter()

	router.Use(nrgorilla.Middleware(newRelicApp))

	countersMdw := wire.InitRequestIdMiddleware(db)
	router.Use(countersMdw.Middleware)

	recoverMdw := recover.NewRecoverMiddleware()
	router.Use(recoverMdw.Middleware)

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

	refreshTokensSubRouter := router.PathPrefix("/refreshtokens").Subrouter()
	refreshTokensSubRouter.Handle("", s.getHandler(authInfra.GetAllRefreshTokensHandler)).Methods(http.MethodGet)
	refreshTokensSubRouter.Handle("", s.getHandler(authInfra.DeleteRefreshTokensHandler)).Methods(http.MethodDelete)
	refreshTokensSubRouter.Use(authMdw.Middleware)
	refreshTokensSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/login", s.getHandler(authInfra.LoginHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(authInfra.RefreshTokenHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(authInfra.RefreshTokenHandler)).Methods(http.MethodDelete)

	pprofSubRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofSubRouter.Handle("/heap", pprof.Handler("heap"))

	logMdw := wire.InitLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	s.addSubscriber(listsInfra.NewListItemCreatedEventSubscriber(s.eventBus, s.listsRepo))
	s.addSubscriber(listsInfra.NewListItemDeletedEventSubscriber(s.eventBus, s.listsRepo))

	s.startSubscribers()

	return &s
}

func (s *server) getHandler(handlerFunc handler.HandlerFunc) handler.Handler {
	return handler.NewHandler(handlerFunc, s.authRepo, s.listsRepo, s.cfgSrv, s.tokenSrv, s.passGen, s.eventBus)
}

func (s *server) addSubscriber(subscriber events.Subscriber) {
	s.subscribers = append(s.subscribers, subscriber)
}

func (s *server) startSubscribers() {
	for _, subscriber := range s.subscribers {
		subscriber.Subscribe()
		go subscriber.Start()
	}
}
