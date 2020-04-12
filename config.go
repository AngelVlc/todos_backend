package main

import (
	"fmt"
	"os"
)

func getDasource() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")

	return fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
}

func getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
