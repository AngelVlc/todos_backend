package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestUsersService(t *testing.T) {
	columns := []string{"id", "name", "password_hash", "is_admin"}
	hasshedPass := "hassedPassword"
	user := "user"

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	mockedCh := MockedCryptoHelper{}

	svc := NewDefaultUsersService(&mockedCh, db)

	t.Run("FindUserByName() should not return a user if it does not exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns))

		u, err := svc.FindUserByName(user)

		assert.Nil(t, u)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindUserByName() should return an error if the query fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnError(fmt.Errorf("some error"))

		u, err := svc.FindUserByName(user)

		assert.Nil(t, u)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user name", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindUserByName() should return the user if it exists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		u, err := svc.FindUserByName(user)

		assert.NotNil(t, u)
		assert.Equal(t, u.ID, int32(5))
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

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

	t.Run("FindUserByID() should not return a user if it does not exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(columns))

		u, err := svc.FindUserByID(1)

		assert.Nil(t, u)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindUserByID() should return an error if the query fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(1).
			WillReturnError(fmt.Errorf("some error"))

		u, err := svc.FindUserByID(1)

		assert.Nil(t, u)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user id", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindUserByID() should return the user if it exists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(5).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		u, err := svc.FindUserByID(5)

		assert.NotNil(t, u)
		assert.Equal(t, u.Name, user)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddUser() should return an error if passwords does not match", func(t *testing.T) {
		dto := dtos.UserDto{
			NewPassword:        "a",
			ConfirmNewPassword: "b",
		}

		_, err := svc.AddUser(&dto)

		appErrors.CheckBadRequestError(t, err, "Passwords don't match", "")
	})

	t.Run("AddUser() should return an error if find user by name fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnError(fmt.Errorf("some error"))

		_, err := svc.AddUser(&dto)

		assert.NotNil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddUser() should return an error if already exists a user with the same name", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		_, err := svc.AddUser(&dto)

		appErrors.CheckBadRequestError(t, err, "A user with the same user name already exists", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddUser() should return an error if generating hassed password fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns))

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(""), fmt.Errorf("some error")).Once()

		_, err := svc.AddUser(&dto)

		appErrors.CheckUnexpectedError(t, err, "Error encrypting password", "some error")

		mockedCh.AssertExpectations(t)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddUser() should return an error if inserting the new user fails", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns))

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`password_hash`,`is_admin`) VALUES (?,?,?)")).
			WithArgs(dto.Name, hasshedPass, false).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUser(&dto)

		appErrors.CheckUnexpectedError(t, err, "Error inserting in the database", "some error")

		mockedCh.AssertExpectations(t)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FindUserByID() should insert the new user", func(t *testing.T) {
		dto := dtos.UserDto{
			Name:               user,
			NewPassword:        "a",
			ConfirmNewPassword: "a",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user).
			WillReturnRows(sqlmock.NewRows(columns))

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		var affected int64
		result := sqlmock.NewResult(12, affected)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`password_hash`,`is_admin`) VALUES (?,?,?)")).
			WithArgs(dto.Name, hasshedPass, false).
			WillReturnResult(result)
		mock.ExpectCommit()

		id, err := svc.AddUser(&dto)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		mockedCh.AssertExpectations(t)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetUsers() should return an error if the query fails", func(t *testing.T) {
		dto := []dtos.GetUserResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,is_admin FROM `users`")).
			WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUsers(&dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting users", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GetUsers() should return the users", func(t *testing.T) {
		dto := []dtos.GetUserResultDto{}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,is_admin FROM `users`")).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user1", "pass1", true).AddRow(12, "user2", "pass2", false))

		err := svc.GetUsers(&dto)

		assert.Equal(t, len(dto), 2)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUser() should return an error if finding the user fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnError(fmt.Errorf("some error"))

		err := svc.RemoveUser(11)

		appErrors.CheckUnexpectedError(t, err, "Error getting user by user id", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUser() should return an error when deleting the admin user", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(5)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, "admin", "", true))

		err := svc.RemoveUser(5)

		appErrors.CheckBadRequestError(t, err, "It is not possible to delete the admin user", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUser() should return an error when deleting a user that does not exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(5)).
			WillReturnRows(sqlmock.NewRows(columns))

		err := svc.RemoveUser(5)

		appErrors.CheckBadRequestError(t, err, "The user does not exist", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUser() should return an error if delete fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, "user", "", true))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE (`users`.`id` = ?)")).
			WithArgs(11).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUser(11)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RemoveUser() should delete the user", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, "user", "", true))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE (`users`.`id` = ?)")).
			WithArgs(11).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUser(11)

		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() should return an error if finding the user fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnError(fmt.Errorf("some error"))

		u, err := svc.UpdateUser(11, &dtos.UserDto{})

		assert.Nil(t, u)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user id", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() should return an error when trying to update the admin user name", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "admin", "", true))

		dto := dtos.UserDto{
			Name:               "anotherName",
			NewPassword:        "a",
			ConfirmNewPassword: "b",
			IsAdmin:            true,
		}

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
		appErrors.CheckBadRequestError(t, err, "It is not possible to change the admin user name", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() should return an error when trying to update the admin is admin field", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "admin", "", true))

		dto := dtos.UserDto{
			Name:               "admin",
			NewPassword:        "a",
			ConfirmNewPassword: "b",
			IsAdmin:            false,
		}

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
		appErrors.CheckBadRequestError(t, err, "It is not possible to change the admin's is admin field", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() should return an error when trying to update the user without changing its password", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user", "hash", true))

		dto := dtos.UserDto{
			Name:    "user",
			IsAdmin: false,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name` = ?, `password_hash` = ?, `is_admin` = ? WHERE `users`.`id` = ?")).
			WithArgs("user", "hash", false, 11).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
		appErrors.CheckUnexpectedError(t, err, "Error updating user", "some error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() should return an error when trying to update the user changing its password but the passwords don't match", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user", "hash", true))

		dto := dtos.UserDto{
			Name:               "user",
			IsAdmin:            false,
			NewPassword:        "a",
			ConfirmNewPassword: "b",
		}

		u, err := svc.UpdateUser(11, &dto)

		assert.Nil(t, u)
		appErrors.CheckBadRequestError(t, err, "Passwords don't match", "")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateUser() to update the user changing its password when the passwords match", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user", "hash", true))

		dto := dtos.UserDto{
			Name:               "user",
			IsAdmin:            false,
			NewPassword:        "new",
			ConfirmNewPassword: "new",
		}

		mockedCh.On("GenerateFromPassword", []byte(dto.NewPassword)).Return([]byte(hasshedPass), nil).Once()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name` = ?, `password_hash` = ?, `is_admin` = ? WHERE `users`.`id` = ?")).
			WithArgs("user", hasshedPass, false, 11).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE `users`.`id` = ? ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user", "", false))

		u, err := svc.UpdateUser(11, &dto)

		assert.NotNil(t, u)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
