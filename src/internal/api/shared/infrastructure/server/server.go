package server

import (
	"net/http"
	"net/http/pprof"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain/passgen"
	authInfra "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure"
	authHandlers "github.com/AngelVlc/todos_backend/src/internal/api/auth/infrastructure/handlers"
	listsDomain "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	listsInfra "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure"
	listsHandlers "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/handlers"
	listSubscribers "github.com/AngelVlc/todos_backend/src/internal/api/lists/infrastructure/subscribers"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/middlewares/recover"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/search"
	"github.com/AngelVlc/todos_backend/src/internal/api/wire"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type server struct {
	http.Handler
	authRepo          authDomain.AuthRepository
	usersRepo         authDomain.UsersRepository
	listsRepo         listsDomain.ListsRepository
	cfgSrv            sharedApp.ConfigurationService
	tokenSrv          authDomain.TokenService
	passGen           passgen.PasswordGenerator
	eventBus          events.EventBus
	subscribers       []events.Subscriber
	newRelicApp       *newrelic.Application
	listsSearchClient search.SearchIndexClient
}

func NewServer(db *gorm.DB, eb events.EventBus, newRelicApp *newrelic.Application) *server {

	s := server{
		authRepo:          wire.InitAuthRepository(db),
		usersRepo:         wire.InitUsersRepository(db),
		listsRepo:         wire.InitListsRepository(db),
		cfgSrv:            wire.InitConfigurationService(),
		tokenSrv:          wire.InitTokenService(),
		passGen:           wire.InitPasswordGenerator(),
		eventBus:          eb,
		subscribers:       []events.Subscriber{},
		newRelicApp:       newRelicApp,
		listsSearchClient: wire.InitSearchIndexClient("lists"),
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
	listsSubRouter.Handle("", s.getHandler(listsHandlers.GetAllListsHandler, nil)).Methods(http.MethodGet)
	listsSubRouter.Handle("", s.getHandler(listsHandlers.CreateListHandler, &listsInfra.ListInput{})).Methods(http.MethodPost)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.GetListHandler, nil)).Methods(http.MethodGet)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.DeleteListHandler, nil)).Methods(http.MethodDelete)
	listsSubRouter.Handle("/{id:[0-9]+}", s.getHandler(listsHandlers.UpdateListHandler, &listsInfra.ListInput{})).Methods(http.MethodPatch)
	listsSubRouter.Handle("/{id:[0-9]+}/move_item", s.getHandler(listsHandlers.MoveListItemHandler, &listsInfra.MoveListItemInput{})).Methods(http.MethodPost)
	listsSubRouter.Use(authMdw.Middleware)

	toolsSubRouter := router.PathPrefix("/tools").Subrouter()
	toolsSubRouter.Handle("/index-lists", s.getHandler(listsHandlers.IndexAllListsHandler, nil)).Methods(http.MethodPost)
	toolsSubRouter.Use(authMdw.Middleware)
	toolsSubRouter.Use(requireAdminMdw.Middleware)

	usersSubRouter := router.PathPrefix("/users").Subrouter()
	usersSubRouter.Handle("", s.getHandler(authHandlers.CreateUserHandler, &authInfra.CreateUserInput{})).Methods(http.MethodPost)
	usersSubRouter.Handle("", s.getHandler(authHandlers.GetAllUsersHandler, nil)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.GetUserHandler, nil)).Methods(http.MethodGet)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.DeleteUserHandler, nil)).Methods(http.MethodDelete)
	usersSubRouter.Handle("/{id:[0-9]+}", s.getHandler(authHandlers.UpdateUserHandler, &authInfra.UpdateUserInput{})).Methods(http.MethodPatch)
	usersSubRouter.Use(authMdw.Middleware)
	usersSubRouter.Use(requireAdminMdw.Middleware)

	refreshTokensSubRouter := router.PathPrefix("/refreshtokens").Subrouter()
	refreshTokensSubRouter.Handle("", s.getHandler(authHandlers.GetAllRefreshTokensHandler, nil)).Methods(http.MethodGet)
	refreshTokensSubRouter.Handle("", s.getHandler(authHandlers.DeleteRefreshTokensHandler, &[]int32{})).Methods(http.MethodDelete)
	refreshTokensSubRouter.Use(authMdw.Middleware)
	refreshTokensSubRouter.Use(requireAdminMdw.Middleware)

	authSubRouter := router.PathPrefix("/auth").Subrouter()
	authSubRouter.Handle("/login", s.getHandler(authHandlers.LoginHandler, &authInfra.LoginInput{})).Methods(http.MethodPost)
	authSubRouter.Handle("/refreshtoken", s.getHandler(authHandlers.RefreshTokenHandler, nil)).Methods(http.MethodPost)
	authSubRouter.Handle("/create_admin", s.getHandler(authHandlers.CreateUserHandler, &authInfra.CreateUserInput{})).Methods(http.MethodPost)

	pprofSubRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofSubRouter.Handle("/heap", pprof.Handler("heap"))

	logMdw := wire.InitLogMiddleware()
	router.Use(logMdw.Middleware)

	s.Handler = router

	s.addSubscriber(listSubscribers.NewListItemsCountProcessor(events.ListCreated, s.eventBus, s.listsRepo, s.newRelicApp))
	s.addSubscriber(listSubscribers.NewListItemsCountProcessor(events.ListUpdated, s.eventBus, s.listsRepo, s.newRelicApp))
	s.addSubscriber(listSubscribers.NewIndexAllListsProcessor(events.IndexAllListsRequested, s.eventBus, s.listsRepo, s.listsSearchClient, s.newRelicApp))
	s.addSubscriber(listSubscribers.NewUpdateSearchIndexDocumentProcessor(events.ListCreated, s.eventBus, s.listsRepo, s.listsSearchClient, s.newRelicApp))
	s.addSubscriber(listSubscribers.NewUpdateSearchIndexDocumentProcessor(events.ListUpdated, s.eventBus, s.listsRepo, s.listsSearchClient, s.newRelicApp))
	s.addSubscriber(listSubscribers.NewRemoveSearchIndexDocumentProcessor(events.ListDeleted, s.eventBus, s.listsSearchClient, s.newRelicApp))

	s.startSubscribers()

	return &s
}

func (s *server) getHandler(handlerFunc handler.HandlerFunc, requestInput interface{}) handler.Handler {
	return handler.NewHandler(handlerFunc, s.authRepo, s.usersRepo, s.listsRepo, s.cfgSrv, s.tokenSrv, s.passGen, s.eventBus, requestInput)
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
