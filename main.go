package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	db, err := gorm.Open("mysql", getDasource())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	checkAdminUser(db)

	fmt.Println("hola caracola")
}
