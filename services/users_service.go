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

// AddUser  adds a user
func (s *UsersService) AddUser(dto *models.UserDto) (int32, error) {
	if dto.NewPassword != dto.ConfirmNewPassword {
		return -1, &appErrors.BadRequestError{Msg: "Passwords don't match", InternalError: nil}
	}

	userExists, err := s.existsUser(dto.Name)
	if err != nil {
		return -1, err
	}

	if userExists {
		return -1, &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	user := dto.ToUser()

	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), 10)
	if err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user.PasswordHash = string(hasshedPass)

	err = s.db.Create(&user).Error
	if err != nil {
		return -1, &appErrors.UnexpectedError{
			Msg:           "Error inserting in the database",
			InternalError: err,
		}
	}

	return user.ID, nil
}

func (s *UsersService) existsUser(userName string) (bool, error) {
	var foundUsers int32
	if err := s.db.Where(models.User{Name: userName}).Table("users").Count(&foundUsers).Error; err != nil {
		return false, &appErrors.UnexpectedError{Msg: "Error checking if user name exists", InternalError: err}
	}
	return foundUsers > 0, nil
}
