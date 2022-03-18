package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AngelVlc/todos_backend/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/server"
	"github.com/AngelVlc/todos_backend/internal/api/wire"
	"github.com/gorilla/handlers"
	"github.com/honeybadger-io/honeybadger-go"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

func main() {
	cfg := wire.InitConfigurationService()

	initHoneyBadger(cfg)
	defer honeybadger.Monitor()

	newRelicApp, err := initNewRelic(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := wire.InitAuthRepository(db)

	go initDeleteExpiredTokensProcess(cfg, authRepo, newRelicApp)

	eb := wire.InitEventBus(map[string]events.DataChannelSlice{})

	server := server.NewServer(db, eb, newRelicApp)

	port := cfg.GetPort()

	address := fmt.Sprintf(":%v", port)

	validCorsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	validCorsOrigins := handlers.AllowedOrigins(cfg.GetCorsAllowedOrigins())
	validCorsMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})
	allowCredentials := handlers.AllowCredentials()

	ctx, cancel := context.WithCancel(context.Background())

	httpServer := &http.Server{
		Addr:         address,
		Handler:      handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(server),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		log.Printf("Listening on port %v ...\n", port)

		if err = http.ListenAndServe(address, handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(server)); err != nil {
			log.Fatalf("could not listen on port %v %v", port, err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Print("os.Interrupt - shutting down...\n")

	go func() {
		<-sigChan
		log.Fatal("os.Kill - terminating...\n")
	}()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	}

	log.Println("gracefully stopped")

	cancel()

	defer os.Exit(0)
}
