package server

import (
	"net/http"
	"net/http/pprof"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	authHandlers "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/handlers"
	listsDomain "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsInfra "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	listsHandlers "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/handlers"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/recover"
	"github.com/AngelVlc/todos_backend/src/internal/api/wire"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type server struct {
	http.Handler
	authRepo    authDomain.AuthRepository
	usersRepo   authDomain.UsersRepository
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
		usersRepo:   wire.InitUsersRepository(db),
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

	router.HandleFunc("/", rootHandler).Methods(http.MethodGet)

	listsSubRouter := router.PathPrefix("/lists").Subrouter()
	listsSubRouter.Handle("", s.getHandler(listsHandlers.GetAllListsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(listsHandlers.CreateListHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.GetListHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.DeleteListHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.UpdateListHandler)).Methods(http.MethodPut)
	listsSubRouter.Handle("/{listId:[0-9]+}/items", s.getHandler(listsHandlers.GetAllListItemsHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{listId:[0-9]+}/items", s.getHandler(listsHandlers.CreateListItemHandler)).Methods(http.MethodPost)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsHandlers.GetListItemHandler)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsHandlers.DeleteListItemHandler)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{listId:[0-9]+}/items/{id:[0-9]+}", s.getHandler(listsHandlers.UpdateListItemHandler)).Methods(http.MethodPut)
	listsSubRouter.Use(authMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(authHandlers.CreateUserHandler)).Methods(http.MethodPost)
	usersSubRouter.Handle("", s.getHandler(authHandlers.GetAllUsersHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.GetUserHandler)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.DeleteUserHandler)).Methods(http.MethodDelete)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.UpdateUserHandler)).Methods(http.MethodPut)
	usersSubRouter.Use(authMdw.Middleware)
	usersSubRouter.Use(requireAdminMdw.Middleware)

	refreshTokensSubRouter := router.PathPrefix("/refreshtokens").Subrouter()
	refreshTokensSubRouter.Handle("", s.getHandler(authHandlers.GetAllRefreshTokensHandler)).Methods(http.MethodGet)
	refreshTokensSubRouter.Handle("", s.getHandler(authHandlers.DeleteRefreshTokensHandler)).Methods(http.MethodDelete)
	refreshTokensSubRouter.Use(authMdw.Middleware)
	refreshTokensSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/login", s.getHandler(authHandlers.LoginHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(authHandlers.RefreshTokenHandler)).Methods(http.MethodPost)
	authSubRouter.Handle("/createadmin", s.getHandler(authHandlers.CreateUserHandler)).Methods(http.MethodPost)

	pprofSubRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofSubRouter.Handle("/heap", pprof.Handler("heap"))

	logMdw := wire.InitLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	s.addSubscriber(listsInfra.NewListItemCreatedEventSubscriber(s.eventBus, s.listsRepo, s.newRelicApp))
	s.addSubscriber(listsInfra.NewListItemDeletedEventSubscriber(s.eventBus, s.listsRepo, s.newRelicApp))

	s.startSubscribers()

	return &s
}

func (s *server) getHandler(handlerFunc handler.HandlerFunc) handler.Handler {
	return handler.NewHandler(handlerFunc, s.authRepo, s.usersRepo, s.listsRepo, s.cfgSrv, s.tokenSrv, s.passGen, s.eventBus)
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
