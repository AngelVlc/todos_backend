package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AngelVlc/todos/dtos"

	"github.com/AngelVlc/todos/services"
	"github.com/AngelVlc/todos/wire"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg := wire.InitConfigurationService()

	db, err := initDb(cfg)
	// db.LogMode(true)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usrSvc := wire.InitUsersService(db)

	foundAdmin, err := usrSvc.FindUserByName("admin")
	if err != nil {
		log.Fatal(err)
	}

	if foundAdmin == nil {
		adminPass := cfg.GetAdminPassword()
		dto := dtos.UserDto{
			Name:               "admin",
			NewPassword:        adminPass,
			ConfirmNewPassword: adminPass,
			IsAdmin:            true,
		}
		_, err = usrSvc.AddUser(&dto)
		if err != nil {
			log.Fatal(err)
		}
	}

	countSvc := wire.InitCountersService(db)
	err = countSvc.CreateCounterIfNotExists("requests")
	if err != nil {
		log.Fatal("error checking requests counter: ", err)
	}

	s := newServer(db)

	port := cfg.GetPort()
	address := fmt.Sprintf(":%v", port)

	validCorsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	validCorsOrigins := handlers.AllowedOrigins(cfg.GetCorsAllowedOrigins())
	validCorsMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})

	log.Printf("Listening on port %v ...\n", port)
	if err = http.ListenAndServe(address, handlers.CORS(validCorsHeaders, validCorsOrigins, validCorsMethods)(s)); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}

func initDb(c services.ConfigurationService) (*gorm.DB, error) {
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
