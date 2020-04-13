package services

import (
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) UsersService {
	return UsersService{db}
}

func (s *UsersService) CreateAdminIfNotExists(password string) error {
	hashedPass, err := s.getPasswordHash(password)
	if err != nil {
		return err
	}

	var user models.User
	s.db.Where(models.User{Name: "admin"}).Attrs(models.User{PasswordHash: hashedPass, IsAdmin: true}).FirstOrCreate(&user)

	return nil
}

func (s *UsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
