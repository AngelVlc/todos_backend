package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
)

type UsersService struct {
	crypto CryptoHelper
	db     *gorm.DB
}

func NewUsersService(crypto CryptoHelper, db *gorm.DB) UsersService {
	return UsersService{crypto, db}
}

func (s *UsersService) FindUserByName(name string) (*models.User, error) {
	foundUser := models.User{}
	err := s.db.Where(models.User{Name: name}).Table("users").First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user name", InternalError: err}
	}

	return &foundUser, nil
}

// CheckIfUserPasswordIsOk returns nil if the password is correct or an error if it isn't
func (s *UsersService) CheckIfUserPasswordIsOk(user *models.User, password string) error {
	return s.crypto.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

// FindUserByID returns a single user from its id
func (s *UsersService) FindUserByID(id int32) (*models.User, error) {
	foundUser := models.User{}
	err := s.db.Where(models.User{ID: id}).Table("users").First(&foundUser).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting user by user id", InternalError: err}
	}

	return &foundUser, nil
}

// AddUser  adds a user
func (s *UsersService) AddUser(dto *dtos.UserDto) (int32, error) {
	if dto.NewPassword != dto.ConfirmNewPassword {
		return -1, &appErrors.BadRequestError{Msg: "Passwords don't match", InternalError: nil}
	}

	foundUser, err := s.FindUserByName(dto.Name)
	if err != nil {
		return -1, err
	}

	if foundUser != nil {
		return -1, &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	user := dto.ToUser()

	hasshedPass, err := s.getPasswordHash(dto.NewPassword)
	if err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user.PasswordHash = hasshedPass

	err = s.db.Create(&user).Error
	if err != nil {
		return -1, &appErrors.UnexpectedError{
			Msg:           "Error inserting in the database",
			InternalError: err,
		}
	}

	return user.ID, nil
}

func (s *UsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := s.crypto.GenerateFromPassword([]byte(p))
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
