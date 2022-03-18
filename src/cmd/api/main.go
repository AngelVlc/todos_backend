package main

import (
	"context"
	"crypto/tls"
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
	"github.com/AngelVlc/todos_backend/pkg/autocerts3cache"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/handlers"
	"github.com/honeybadger-io/honeybadger-go"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"golang.org/x/crypto/acme/autocert"
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

	validCorsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	validCorsOrigins := handlers.AllowedOrigins(cfg.GetCorsAllowedOrigins())
	validCorsMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})
	allowCredentials := handlers.AllowCredentials()

	ctx, cancel := context.WithCancel(context.Background())

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.GetPort()),
		Handler:      handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(server),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	var certManager *autocert.Manager

	if cfg.InProduction() {
		ctx := context.TODO()
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("error loading the default config: %v", err)
		}

		awsS3Api := autocerts3cache.NewAwsS3Api(awsCfg)
		s3Cache := autocerts3cache.NewS3Cache(cfg.GetBucketName(), awsS3Api)

		certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.GetDns()),
			Cache:      s3Cache,
		}

		tlsConfig := &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
		tlsConfig.NextProtos = append([]string{"h2", "http/1.1", "acme-tls/1"}, tlsConfig.NextProtos...)

		httpServer.TLSConfig = tlsConfig
		httpServer.Addr = ":443"
	}

	go func() {
		log.Printf("Starting listener on port %v\n", httpServer.Addr)

		if cfg.InProduction() {
			err = httpServer.ListenAndServeTLS("", "")
		} else {
			err = httpServer.ListenAndServe()
		}

		if err != nil {
			log.Fatalf("could not listen on port %v %v", httpServer.Addr, err)
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
