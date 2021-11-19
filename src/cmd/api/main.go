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

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/server"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos/internal/api/wire"
	"github.com/gorilla/handlers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := wire.InitConfigurationService()

	db, err := initDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	createAdminUserIfNotExists(cfg, db)

	authRepo := wire.InitAuthRepository(db)
	initDeleteExpiredTokensProcess(authRepo)

	countersRepo := wire.InitCountersRepository(db)
	countSvc := sharedApp.NewInitRequestsCounterService(countersRepo)
	err = countSvc.InitRequestsCounter()
	if err != nil {
		log.Fatal("error checking requests counter: ", err)
	}

	eb := wire.InitEventBus(map[string]events.DataChannelSlice{})

	s := server.NewServer(db, eb)

	port := cfg.GetPort()

	address := fmt.Sprintf(":%v", port)

	validCorsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	validCorsOrigins := handlers.AllowedOrigins(cfg.GetCorsAllowedOrigins())
	validCorsMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})
	allowCredentials := handlers.AllowCredentials()

	ctx, cancel := context.WithCancel(context.Background())

	httpServer := &http.Server{
		Addr:         address,
		Handler:      handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(s),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		log.Printf("Listening on port %v ...\n", port)
		if err = http.ListenAndServe(address, handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(s)); err != nil {
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
	} else {
		log.Println("gracefully stopped")
	}

	cancel()

	defer os.Exit(0)
}

func initDb(c sharedApp.ConfigurationService) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.GetDatasource()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createAdminUserIfNotExists(cfg sharedApp.ConfigurationService, db *gorm.DB) {
	repo := wire.InitAuthRepository(db)

	userName := authDomain.UserName("admin")
	adminExists, err := repo.ExistsUser(userName)
	if err != nil {
		log.Fatal(err)

		return
	}

	if !adminExists {
		passGen := wire.InitPasswordGenerator()
		hassedPass, _ := passGen.GenerateFromPassword(cfg.GetAdminPassword())

		user := authDomain.User{
			Name:         "admin",
			PasswordHash: hassedPass,
			IsAdmin:      true,
		}
		err = repo.CreateUser(&user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initDeleteExpiredTokensProcess(authRepo authDomain.AuthRepository) {
	ticker := time.NewTicker(30 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				authRepo.DeleteExpiredRefreshTokens(t)
			}
		}
	}()
}
