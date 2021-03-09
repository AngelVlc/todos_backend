package services

import (
	"strings"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/models"
	"github.com/AngelVlc/todos/internal/api/repositories"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/stretchr/testify/mock"
)

type UsersService interface {
	FindUserByID(id int32) (*models.User, error)
	UpdateUser(id int32, dto *dtos.UserDto) error
}

type MockedUsersService struct {
	mock.Mock
}

func NewMockedUsersService() *MockedUsersService {
	return &MockedUsersService{}
}

func (m *MockedUsersService) FindUserByID(id int32) (*models.User, error) {
	args := m.Called(id)
	got := args.Get(0)
	if got == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
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

// FindUserByID returns a single user from its id
func (s *DefaultUsersService) FindUserByID(id int32) (*models.User, error) {
	return s.usersRepo.FindByID(id)
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
