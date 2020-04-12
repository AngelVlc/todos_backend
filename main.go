package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", getDasource())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hola caracola")
}

func getDasource() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user:= os.Getenv("MYSQL_USER")
	pass:= os.Getenv("MYSQL_PASSWORD")
	dbname:= os.Getenv("MYSQL_DATABASE")

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, dbname)
}
