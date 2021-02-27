package services

import (
	"fmt"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedCryptoHelper struct {
	mock.Mock
}

func (m *MockedCryptoHelper) GenerateFromPassword(password []byte) ([]byte, error) {
	args := m.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockedCryptoHelper) CompareHashAndPassword(hashedPassword, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

var (
	columns     = []string{"id", "name", "password_hash", "is_admin"}
	hasshedPass = "hassedPassword"
	user        = "user"
)

func TestGetUsers(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}

	svc := NewDefaultUsersService(nil, &mockedUsersRepo)

	t.Run("should return an error if repository GetAll fails", func(t *testing.T) {
		mockedUsersRepo.On("GetAll").Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.GetUsers()

		assert.Nil(t, dto)
		assert.Error(t, err)

		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return the users", func(t *testing.T) {
		found := []*models.User{
			{ID: 2, Name: "user1", IsAdmin: true},
			{ID: 5, Name: "user2", IsAdmin: false},
		}

		mockedUsersRepo.On("GetAll").Return(found, nil)

		res, err := svc.GetUsers()

		assert.Nil(t, err)
		require.Equal(t, len(res), 2)
		assert.Equal(t, int32(2), res[0].ID)
		assert.Equal(t, "user1", res[0].Name)
		assert.True(t, res[0].IsAdmin)
		assert.Equal(t, int32(5), res[1].ID)
		assert.Equal(t, "user2", res[1].Name)
		assert.False(t, res[1].IsAdmin)

		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}

	svc := NewDefaultUsersService(nil, &mockedUsersRepo)

	userID := int32(11)
	foundUser := models.User{
		ID:      userID,
		Name:    "userName",
		IsAdmin: true,
	}

	t.Run("should return an error if repository FindByID fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", userID).Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.FindUserByID(11)

		assert.Nil(t, dto)
		assert.Error(t, err)

		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return the found user if repository FindByID doesn't fail", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", userID).Return(&foundUser, nil).Once()

		dto, err := svc.FindUserByID(11)

		require.NotNil(t, dto)
		require.IsType(t, &models.User{}, dto)
		assert.Nil(t, err)

		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestFindByName(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}

	svc := NewDefaultUsersService(nil, &mockedUsersRepo)

	userID := int32(11)
	foundUser := models.User{
		ID:      userID,
		Name:    "userName",
		IsAdmin: true,
	}

	t.Run("should return an error if repository FindByName fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByName", "userName").Return(nil, fmt.Errorf("some error")).Once()

		dto, err := svc.FindUserByName("userName")

		assert.Nil(t, dto)
		assert.Error(t, err)

		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return the found user if repository FindByName doesn't fail", func(t *testing.T) {
		mockedUsersRepo.On("FindByName", "userName").Return(&foundUser, nil).Once()

		dto, err := svc.FindUserByName("userName")

		require.NotNil(t, dto)
		require.IsType(t, &models.User{}, dto)
		assert.Nil(t, err)

		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestAddUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo)

	t.Run("should return an error if passwords does not match", func(t *testing.T) {
		dto := dtos.UserDto{
			NewPassword:        "a",
			ConfirmNewPassword: "b",
		}

		_, err := svc.AddUser(&dto)

		appErrors.CheckBadRequestError(t, err, "Passwords don't match", "")
	})

	t.Run("should return an error if find user by name fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mockedUsersRepo.On("FindByName", user).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUser(&dto)

		assert.Error(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error if already exists a user with the same name", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mockedUsersRepo.On("FindByName", user).Return(&models.User{ID: 11, Name: user}, nil).Once()

		_, err := svc.AddUser(&dto)

		appErrors.CheckBadRequestError(t, err, "A user with the same user name already exists", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error if generating hassed password fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mockedUsersRepo.On("FindByName", user).Return(nil, nil).Once()

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(""), fmt.Errorf("some error")).Once()

		_, err := svc.AddUser(&dto)

		appErrors.CheckUnexpectedError(t, err, "Error encrypting password", "some error")

		mockedCh.AssertExpectations(t)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error if creating the new user fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mockedUsersRepo.On("FindByName", user).Return(nil, nil).Once()

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		user := models.User{}
		user.FromDto(&dto)
		user.PasswordHash = hasshedPass

		mockedUsersRepo.On("Create", &user).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUser(&dto)

		assert.Error(t, err)
		mockedCh.AssertExpectations(t)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should create the new user", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mockedUsersRepo.On("FindByName", user).Return(nil, nil).Once()

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		user := models.User{}
		user.FromDto(&dto)
		user.PasswordHash = hasshedPass

		mockedUsersRepo.On("Create", &user).Return(int32(12), nil).Once()

		id, err := svc.AddUser(&dto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		mockedCh.AssertExpectations(t)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestRemoveUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo)

	t.Run("should return an error if finding the user fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.RemoveUser(11)

		require.Error(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when deleting the admin user", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "admin"}, nil).Once()

		err := svc.RemoveUser(11)

		appErrors.CheckBadRequestError(t, err, "It is not possible to delete the admin user", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when deleting a user that does not exist", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(nil, nil).Once()

		err := svc.RemoveUser(11)

		appErrors.CheckBadRequestError(t, err, "The user does not exist", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11}, nil).Once()
		mockedUsersRepo.On("Delete", int32(11)).Return(fmt.Errorf("some error")).Once()

		err := svc.RemoveUser(11)

		assert.Error(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should delete the user", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11}, nil).Once()
		mockedUsersRepo.On("Delete", int32(11)).Return(nil).Once()

		err := svc.RemoveUser(11)

		assert.Nil(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo)

	t.Run("should return an error if finding the user fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(nil, fmt.Errorf("some error")).Once()

		err := svc.UpdateUser(11, &dtos.UserDto{})

		require.Error(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error if the user fails does not exist", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(nil, nil).Once()

		err := svc.UpdateUser(11, &dtos.UserDto{})

		appErrors.CheckBadRequestError(t, err, "The user does not exist", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when trying to update the admin user name", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "admin", IsAdmin: true}, nil).Once()

		dto := dtos.UserDto{
			Name:               "anotherName",
			NewPassword:        "a",
			ConfirmNewPassword: "b",
			IsAdmin:            true,
		}

		err := svc.UpdateUser(11, &dto)

		appErrors.CheckBadRequestError(t, err, "It is not possible to change the admin user name", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when trying to update the admin is admin field", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "admin", IsAdmin: true}, nil).Once()

		dto := dtos.UserDto{
			Name:               "admin",
			NewPassword:        "a",
			ConfirmNewPassword: "b",
			IsAdmin:            false,
		}

		err := svc.UpdateUser(11, &dto)

		appErrors.CheckBadRequestError(t, err, "It is not possible to change the admin's is admin field", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when trying to update the user without changing its password and the update fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:    "user",
			IsAdmin: false,
		}

		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "user", PasswordHash: hasshedPass}, nil).Once()
		mockedUsersRepo.On("Update", &models.User{ID: 11, Name: "user", PasswordHash: hasshedPass}).Return(fmt.Errorf("some error")).Once()

		err := svc.UpdateUser(11, &dto)

		require.Error(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should return an error when trying to update the user changing its password but the passwords don't match", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "user", PasswordHash: "hash"}, nil).Once()

		dto := dtos.UserDto{
			Name:               "user",
			IsAdmin:            false,
			NewPassword:        "a",
			ConfirmNewPassword: "b",
		}

		err := svc.UpdateUser(11, &dto)

		appErrors.CheckBadRequestError(t, err, "Passwords don't match", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should update the user changing its password when the passwords match", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "user", PasswordHash: "hash"}, nil).Once()

		dto := dtos.UserDto{
			Name:               "user",
			IsAdmin:            false,
			NewPassword:        "new",
			ConfirmNewPassword: "new",
		}

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		mockedUsersRepo.On("Update", &models.User{ID: 11, Name: "user", PasswordHash: hasshedPass}).Return(nil).Once()

		err := svc.UpdateUser(11, &dto)

		assert.Nil(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestCheckIfUserPasswordIsOk(t *testing.T) {
	mockedCh := MockedCryptoHelper{}
	mockedUsersRepo := repositories.MockedUsersRepository{}
	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo)

	t.Run("should return nil if the password is ok", func(t *testing.T) {
		user := models.User{
			Name:         "wadus",
			PasswordHash: "hash",
		}

		mockedCh.On("CompareHashAndPassword", []byte(user.PasswordHash), []byte("pass")).Return(nil).Once()

		err := svc.CheckIfUserPasswordIsOk(&user, "pass")

		assert.Nil(t, err)
		mockedCh.AssertExpectations(t)
	})

	t.Run("should return an error if the password is not ok", func(t *testing.T) {
		user := models.User{
			Name:         "wadus",
			PasswordHash: "hash",
		}

		mockedCh.On("CompareHashAndPassword", []byte(user.PasswordHash), []byte("pass")).Return(fmt.Errorf("some error")).Once()

		err := svc.CheckIfUserPasswordIsOk(&user, "pass")

		assert.Error(t, err)
		appErrors.CheckErrorMsg(t, err, "some error")
		mockedCh.AssertExpectations(t)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
