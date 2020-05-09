package services

import (
	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type UsersService interface {
	FindUserByName(name string) (*models.User, error)
	CheckIfUserPasswordIsOk(user *models.User, password string) error
	FindUserByID(id int32) (*models.User, error)
	AddUser(dto *dtos.UserDto) (int32, error)
	GetUsers(r *[]dtos.GetUsersResultDto) error
}

type MockedUsersService struct {
	mock.Mock
}

func NewMockedUsersService() *MockedUsersService {
	return &MockedUsersService{}
}

func (m *MockedUsersService) FindUserByName(name string) (*models.User, error) {
	args := m.Called(name)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockedUsersService) CheckIfUserPasswordIsOk(user *models.User, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func (m *MockedUsersService) FindUserByID(id int32) (*models.User, error) {
	args := m.Called(id)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockedUsersService) AddUser(dto *dtos.UserDto) (int32, error) {
	args := m.Called(dto)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockedUsersService) GetUsers(r *[]dtos.GetUsersResultDto) error {
	args := m.Called(r)
	return args.Error(0)
}

type DefaultUsersService struct {
	crypto CryptoHelper
	db     *gorm.DB
}

func NewDefaultUsersService(crypto CryptoHelper, db *gorm.DB) *DefaultUsersService {
	return &DefaultUsersService{crypto, db}
}

func (s *DefaultUsersService) FindUserByName(name string) (*models.User, error) {
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
func (s *DefaultUsersService) CheckIfUserPasswordIsOk(user *models.User, password string) error {
	return s.crypto.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

// FindUserByID returns a single user from its id
func (s *DefaultUsersService) FindUserByID(id int32) (*models.User, error) {
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
func (s *DefaultUsersService) AddUser(dto *dtos.UserDto) (int32, error) {
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

func (s *DefaultUsersService) GetUsers(r *[]dtos.GetUsersResultDto) error {
	if err := s.db.Find(&r).Error; err != nil {
		return &appErrors.UnexpectedError{Msg: "Error getting users", InternalError: err}
	}
	return nil
}

func (s *DefaultUsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := s.crypto.GenerateFromPassword([]byte(p))
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
