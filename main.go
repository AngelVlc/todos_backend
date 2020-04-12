package main

import (
	"fmt"
	"log"
	"database/sql"
  _"github.com/go-sql-driver/mysql"
)

func main() {
	_, err := sql.Open("mysql", "root:pass@tcp(mysql:3306)/todos")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("hola caracola")
}
