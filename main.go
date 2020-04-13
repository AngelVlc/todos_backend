package main

import (
	"fmt"
	"github.com/AngelVlc/todos/services"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	cfg := services.NewConfigurationService()

	db, err := initDb(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usrSvc := initUsersService(db, &cfg)
	
	err = usrSvc.CreateAdminIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hola caracola")
}

func initDb(c *services.ConfigurationService) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", c.GetDasource())
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
