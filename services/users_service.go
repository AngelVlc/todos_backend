package services

import (
	"strings"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/stretchr/testify/mock"
)

type UsersService interface {
	FindUserByName(name string) (*models.User, error)
	CheckIfUserPasswordIsOk(user *models.User, password string) error
	FindUserByID(id int32) (*models.User, error)
	AddUser(dto *dtos.UserDto) (int32, error)
	GetUsers() ([]*dtos.UserResponseDto, error)
	RemoveUser(id int32) error
	UpdateUser(id int32, dto *dtos.UserDto) error
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

func (m *MockedUsersService) GetUsers() ([]*dtos.UserResponseDto, error) {
	args := m.Called()
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dtos.UserResponseDto), args.Error(1)
}

func (m *MockedUsersService) RemoveUser(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedUsersService) UpdateUser(id int32, dto *dtos.UserDto) error {
	args := m.Called(id, dto)
	return args.Error(0)
}

type DefaultUsersService struct {
	crypto    CryptoHelper
	usersRepo repositories.UsersRepository
}

func NewDefaultUsersService(crypto CryptoHelper, usersRepo repositories.UsersRepository) *DefaultUsersService {
	return &DefaultUsersService{crypto, usersRepo}
}

func (s *DefaultUsersService) FindUserByName(name string) (*models.User, error) {
	return s.usersRepo.FindByName(name)
}

// CheckIfUserPasswordIsOk returns nil if the password is correct or an error if it isn't
func (s *DefaultUsersService) CheckIfUserPasswordIsOk(user *models.User, password string) error {
	return s.crypto.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

// FindUserByID returns a single user from its id
func (s *DefaultUsersService) FindUserByID(id int32) (*models.User, error) {
	return s.usersRepo.FindByID(id)
}

// AddUser  adds a user
func (s *DefaultUsersService) AddUser(dto *dtos.UserDto) (int32, error) {
	if dto.NewPassword != dto.ConfirmNewPassword {
		return -1, &appErrors.BadRequestError{Msg: "Passwords don't match", InternalError: nil}
	}

	foundUser, err := s.usersRepo.FindByName(dto.Name)
	if err != nil {
		return -1, err
	}

	if foundUser != nil {
		return -1, &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	user := models.User{}
	user.FromDto(dto)

	hasshedPass, err := s.getPasswordHash(dto.NewPassword)
	if err != nil {
		return -1, &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user.PasswordHash = hasshedPass

	id, err := s.usersRepo.Create(&user)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *DefaultUsersService) GetUsers() ([]*dtos.UserResponseDto, error) {
	found, err := s.usersRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]*dtos.UserResponseDto, len(found))

	for i, v := range found {
		res[i] = v.ToResponseDto()
	}

	return res, nil
}

func (s *DefaultUsersService) RemoveUser(id int32) error {
	foundUser, err := s.usersRepo.FindByID(id)
	if err != nil {
		return err
	}

	if foundUser == nil {
		return &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	if strings.ToLower(foundUser.Name) == "admin" {
		return &appErrors.BadRequestError{Msg: "It is not possible to delete the admin user"}
	}

	return s.usersRepo.Delete(id)
}

func (s *DefaultUsersService) UpdateUser(id int32, dto *dtos.UserDto) error {
	foundUser, err := s.usersRepo.FindByID(id)
	if err != nil {
		return err
	}

	if foundUser == nil {
		return &appErrors.BadRequestError{Msg: "The user does not exist"}
	}

	if strings.ToLower(foundUser.Name) == "admin" {
		if dto.Name != "admin" {
			return &appErrors.BadRequestError{Msg: "It is not possible to change the admin user name"}
		}

		if !dto.IsAdmin {
			return &appErrors.BadRequestError{Msg: "It is not possible to change the admin's is admin field"}
		}
	}

	foundUser.FromDto(dto)

	if len(dto.NewPassword) > 0 {
		if dto.NewPassword != dto.ConfirmNewPassword {
			return &appErrors.BadRequestError{Msg: "Passwords don't match", InternalError: nil}
		}

		hasshedPass, err := s.getPasswordHash(dto.NewPassword)
		if err != nil {
			return &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
		}

		foundUser.PasswordHash = hasshedPass
	}

	return s.usersRepo.Update(foundUser)
}

func (s *DefaultUsersService) getPasswordHash(p string) (string, error) {
	hasshedPass, err := s.crypto.GenerateFromPassword([]byte(p))
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
