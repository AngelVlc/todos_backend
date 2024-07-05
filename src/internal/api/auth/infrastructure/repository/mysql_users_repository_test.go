//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userColumns = []string{"id", "name", "passwordHash", "isAdmin"}
)

func TestMySqlUsersRepository_FindUser_WhenTheQueryFails(t *testing.T) {
	userID := int32(1)
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? LIMIT 1")).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some error"))

	repo := NewMySqlUsersRepository(db)

	res, err := repo.FindUser(context.Background(), domain.UserRecord{ID: userID})

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_FindUser_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`name` = ? LIMIT 1")).
		WithArgs("userName").
		WillReturnRows(sqlmock.NewRows(userColumns).AddRow(1, "userName", "hash", true))

	repo := NewMySqlUsersRepository(db)

	res, err := repo.FindUser(context.Background(), domain.UserRecord{Name: "userName"})

	require.NotNil(t, res)
	assert.Equal(t, "userName", res.Name)
	assert.True(t, res.IsAdmin)
	assert.Equal(t, int32(1), res.ID)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_ExistsUser_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`name` = ?")).
		WithArgs("userName").
		WillReturnError(fmt.Errorf("some error"))

	repo := NewMySqlUsersRepository(db)
	user := domain.UserRecord{Name: "userName"}

	res, err := repo.ExistsUser(context.Background(), user)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepository_ExistsUser_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`name` = ?")).
		WithArgs("userName").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	repo := NewMySqlUsersRepository(db)
	user := domain.UserRecord{Name: "userName"}

	res, err := repo.ExistsUser(context.Background(), user)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_GetAll_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnError(fmt.Errorf("some error"))

	repo := NewMySqlUsersRepository(db)

	res, err := repo.GetAll(context.Background())

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_GetAll_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WillReturnRows(sqlmock.NewRows(userColumns).
			AddRow(11, "user1", "pass1", true).
			AddRow(12, "user2", "pass2", false))

	repo := NewMySqlUsersRepository(db)

	res, err := repo.GetAll(context.Background())

	assert.Nil(t, err)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	nvo, _ := domain.NewUserNameValueObject("user1")
	assert.Equal(t, nvo, res[0].Name)
	assert.Equal(t, "pass1", res[0].PasswordHash)
	assert.True(t, res[0].IsAdmin)
	assert.Equal(t, int32(12), res[1].ID)
	nvo, _ = domain.NewUserNameValueObject("user2")
	assert.Equal(t, nvo, res[1].Name)
	assert.Equal(t, "pass2", res[1].PasswordHash)
	assert.False(t, res[1].IsAdmin)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Create_WhenItFails(t *testing.T) {
	user := domain.UserRecord{Name: "userName", PasswordHash: "hash", IsAdmin: false}
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`) VALUES (?,?,?)")).
		WithArgs(user.Name, user.PasswordHash, user.IsAdmin).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	repo := NewMySqlUsersRepository(db)

	err := repo.Create(context.Background(), &user)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Create_WhenItDoesNotFail(t *testing.T) {
	user := domain.UserRecord{Name: "userName", PasswordHash: "hash", IsAdmin: false}
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`) VALUES (?,?,?)")).
		WithArgs(user.Name, user.PasswordHash, user.IsAdmin).
		WillReturnResult(sqlmock.NewResult(12, 1))
	mock.ExpectCommit()

	repo := NewMySqlUsersRepository(db)

	err := repo.Create(context.Background(), &user)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Delete_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	repo := NewMySqlUsersRepository(db)
	userID := int32(1)

	err := repo.Delete(context.Background(), domain.UserRecord{ID: userID})

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Delete_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewMySqlUsersRepository(db)
	userID := int32(1)

	err := repo.Delete(context.Background(), domain.UserRecord{ID: userID})

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Update_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`passwordHash`=?,`isAdmin`=? WHERE `id` = ?")).
		WithArgs("userName", "hash", false, 11).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	repo := NewMySqlUsersRepository(db)
	user := domain.UserRecord{ID: 11, Name: "userName", PasswordHash: "hash", IsAdmin: false}

	err := repo.Update(context.Background(), &user)

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepository_Update_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`passwordHash`=?,`isAdmin`=? WHERE `id` = ?")).
		WithArgs("userName", "hash", false, 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`,`id`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`passwordHash`=VALUES(`passwordHash`),`isAdmin`=VALUES(`isAdmin`)")).
		WithArgs("userName", "hash", false, 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewMySqlUsersRepository(db)
	user := domain.UserRecord{ID: 11, Name: "userName", PasswordHash: "hash", IsAdmin: false}

	err := repo.Update(context.Background(), &user)

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}
