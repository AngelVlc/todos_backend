package main

import (
	"fmt"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
)

func checkAdminUser(db *gorm.DB) error {
	var user models.User
	adminPass := getAdminPassword()
	hashedPass, err := getPasswordHash(adminPass)
	if err != nil {
		return err
	}

	db.Where(models.User{Name: "admin"}).Attrs(models.User{PasswordHash: hashedPass, IsAdmin: true}).FirstOrCreate(&user)

	fmt.Println("###", user, "###")

	return nil
}
