package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/repositories"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
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

func TestUsersServiceFindByID(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}

	svc := NewDefaultUsersService(nil, &mockedUsersRepo, nil)

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

func TestUsersServiceFindByName(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}

	svc := NewDefaultUsersService(nil, &mockedUsersRepo, nil)

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

func TestUsersServiceAddUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo, nil)

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

		assert.NotNil(t, err)
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

	t.Run("should return an error if inserting the new user fails", func(t *testing.T) {
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

		mockedUsersRepo.On("Insert", &user).Return(nil, fmt.Errorf("some error")).Once()

		_, err := svc.AddUser(&dto)

		assert.NotNil(t, err)
		mockedCh.AssertExpectations(t)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should insert the new user", func(t *testing.T) {
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

		mockedUsersRepo.On("Insert", &user).Return(int32(12), nil).Once()

		id, err := svc.AddUser(&dto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		mockedCh.AssertExpectations(t)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestUsersServiceRemoveUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo, nil)

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
		mockedUsersRepo.On("Remove", int32(11)).Return(fmt.Errorf("some error")).Once()

		err := svc.RemoveUser(11)

		assert.NotNil(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("should delete the user", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11}, nil).Once()
		mockedUsersRepo.On("Remove", int32(11)).Return(nil).Once()

		err := svc.RemoveUser(11)

		assert.Nil(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestUsersServiceUpdateUser(t *testing.T) {
	mockedUsersRepo := repositories.MockedUsersRepository{}
	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo, nil)

	t.Run("should return an error if finding the user fails", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(nil, fmt.Errorf("some error")).Once()

		u, err := svc.UpdateUser(11, &dtos.UserDto{})

		assert.Nil(t, u)
		require.Error(t, err)
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

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
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

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
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

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
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

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
		appErrors.CheckBadRequestError(t, err, "Passwords don't match", "")
		mockedUsersRepo.AssertExpectations(t)
	})

	t.Run("to update the user changing its password when the passwords match", func(t *testing.T) {
		mockedUsersRepo.On("FindByID", int32(11)).Return(&models.User{ID: 11, Name: "user", PasswordHash: "hash"}, nil).Once()

		dto := dtos.UserDto{
			Name:               "user",
			IsAdmin:            false,
			NewPassword:        "new",
			ConfirmNewPassword: "new",
		}

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		mockedUsersRepo.On("Update", &models.User{ID: 11, Name: "user", PasswordHash: hasshedPass}).Return(nil).Once()

		u, err := svc.UpdateUser(11, &dto)

		assert.NotNil(t, u)
		assert.Nil(t, err)
		mockedUsersRepo.AssertExpectations(t)
	})
}

func TestUsersService(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	mockedCh := MockedCryptoHelper{}
	mockedUsersRepo := repositories.MockedUsersRepository{}
	svc := NewDefaultUsersService(&mockedCh, &mockedUsersRepo, db)

	t.Run("CheckIfUserPasswordIsOk() should return nil if the password is ok", func(t *testing.T) {
		user := models.User{
			Name:         "wadus",
			PasswordHash: "hash",
		}

		mockedCh.On("CompareHashAndPassword", []byte(user.PasswordHash), []byte("pass")).Return(nil).Once()

		err := svc.CheckIfUserPasswordIsOk(&user, "pass")

		assert.Nil(t, err)

		mockedCh.AssertExpectations(t)
	})

	t.Run("CheckIfUserPasswordIsOk() should return an error if the password is not ok", func(t *testing.T) {
		user := models.User{
			Name:         "wadus",
			PasswordHash: "hash",
		}

		mockedCh.On("CompareHashAndPassword", []byte(user.PasswordHash), []byte("pass")).Return(fmt.Errorf("some error")).Once()

		err := svc.CheckIfUserPasswordIsOk(&user, "pass")

		assert.NotNil(t, err)
		appErrors.CheckErrorMsg(t, err, "some error")

		mockedCh.AssertExpectations(t)
	})

	expectedGetUsersQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,is_admin FROM `users`"))
	}

	t.Run("GetUsers() should return an error if the query fails", func(t *testing.T) {
		dto := []dtos.GetUserResultDto{}

		expectedGetUsersQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUsers(&dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting users", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetUsers() should return the users", func(t *testing.T) {
		dto := []dtos.GetUserResultDto{}

		expectedGetUsersQuery().WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user1", "pass1", true).AddRow(12, "user2", "pass2", false))

		err := svc.GetUsers(&dto)

		assert.Equal(t, len(dto), 2)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
