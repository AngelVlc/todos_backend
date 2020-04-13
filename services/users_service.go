package services

import (
	appErrors "github.com/AngelVlc/todos/errors"
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

// CheckIfUserPasswordIsOk returns nil if the password is correct or an error if it isn't
func (s *UsersService) CheckIfUserPasswordIsOk(userName string, password string) (*models.User, error) {
	foundUser := s.getUserByUserName(userName)

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist", InternalError: nil}
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password))
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid password", InternalError: nil}
	}

	return foundUser, nil
}

// GetUserByID returns a single user from its id
func (s *UsersService) GetUserByID(id int32) *models.User {
	foundUser := models.User{}

	s.db.Where(models.User{ID: id}).First(&foundUser)

	return &foundUser
}

func (s *UsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}

func (s *UsersService) getUserByUserName(userName string) *models.User {
	foundUser := models.User{}

	s.db.Where(models.User{Name: "admin"}).First(&foundUser)

	return &foundUser
}
