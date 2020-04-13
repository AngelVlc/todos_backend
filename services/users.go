package services

import (
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	db     *gorm.DB
	config *ConfigurationService
}

func NewUsersService(db *gorm.DB, config *ConfigurationService) UsersService {
	return UsersService{db, config}
}

func (u *UsersService) CreateAdminIfNotExists() error {
	adminPass := u.config.GetAdminPassword()
	hashedPass, err := u.getPasswordHash(adminPass)
	if err != nil {
		return err
	}

	var user models.User
	u.db.Where(models.User{Name: "admin"}).Attrs(models.User{PasswordHash: hashedPass, IsAdmin: true}).FirstOrCreate(&user)

	return nil
}

func (u *UsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
