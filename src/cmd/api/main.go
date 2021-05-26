package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/AngelVlc/todos/internal/api/server"
	sharedApp "github.com/AngelVlc/todos/internal/api/shared/application"
	"github.com/AngelVlc/todos/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos/internal/api/wire"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg := wire.InitConfigurationService()

	db, err := initDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//	db.LogMode(true)

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

	log.Printf("Listening on port %v ...\n", port)
	if err = http.ListenAndServe(address, handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods, allowCredentials)(s)); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}

func initDb(c sharedApp.ConfigurationService) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", c.GetDatasource())
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createAdminUserIfNotExists(cfg sharedApp.ConfigurationService, db *gorm.DB) {
	repo := wire.InitAuthRepository(db)

	userName := authDomain.UserName("admin")
	foundAdmin, err := repo.FindUserByName(userName)
	if err != nil {
		log.Fatal(err)
	}

	if foundAdmin == nil {
		passGen := wire.InitPasswordGenerator()
		hassedPass, err := passGen.GenerateFromPassword(cfg.GetAdminPassword())

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
