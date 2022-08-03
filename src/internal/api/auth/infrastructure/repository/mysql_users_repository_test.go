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

func TestMySqlUsersRepositoryFindUser_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? LIMIT 1")).
			WithArgs(userID)
	}

	expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.FindUser(context.Background(), &domain.User{ID: userID})

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryFindUser_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	userName := domain.UserName("userName")

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`name` = ? LIMIT 1")).
			WithArgs(userName)
	}

	expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(int32(1), userName, "hash", true))

	res, err := repo.FindUser(context.Background(), &domain.User{Name: userName})

	require.NotNil(t, res)
	assert.Equal(t, userName, res.Name)
	assert.True(t, res.IsAdmin)
	assert.Equal(t, int32(1), res.ID)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepositoryExistsUser_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	user := &domain.User{Name: domain.UserName("userName")}

	expectedExistsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`name` = ?")).
			WithArgs("userName")
	}

	expectedExistsQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.ExistsUser(context.Background(), user)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlAuthRepositoryExistsUser_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	user := &domain.User{Name: domain.UserName("userName")}

	expectedExistsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`name` = ?")).
			WithArgs("userName")
	}

	expectedExistsQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	res, err := repo.ExistsUser(context.Background(), user)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryGetAll_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	expectedQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`"))
	}

	expectedQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAll(context.Background())

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryGetAll_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	expectedQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`"))
	}

	expectedQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(11, "user1", "pass1", true).AddRow(12, "user2", "pass2", false))

	res, err := repo.GetAll(context.Background())

	assert.Nil(t, err)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, domain.UserName("user1"), res[0].Name)
	assert.Equal(t, "pass1", res[0].PasswordHash)
	assert.True(t, res[0].IsAdmin)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, domain.UserName("user2"), res[1].Name)
	assert.Equal(t, "pass2", res[1].PasswordHash)
	assert.False(t, res[1].IsAdmin)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryCreate_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlUsersRepository(db)

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`) VALUES (?,?,?)")).
			WithArgs(user.Name, user.PasswordHash, user.IsAdmin)
	}

	mock.ExpectBegin()
	expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.Create(context.Background(), &user)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryCreate_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlUsersRepository(db)

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`) VALUES (?,?,?)")).
			WithArgs(user.Name, user.PasswordHash, user.IsAdmin)
	}

	result := sqlmock.NewResult(12, 1)

	mock.ExpectBegin()
	expectedInsertExec().WillReturnResult(result)
	mock.ExpectCommit()

	err := repo.Create(context.Background(), &user)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryDelete_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
			WithArgs(1)
	}

	userID := int32(1)

	mock.ExpectBegin()
	expectedDeleteExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.Delete(context.Background(), &domain.User{ID: userID})

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)

}

func TestMySqlUsersRepositoryDelete_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
			WithArgs(1)
	}

	userID := int32(1)

	mock.ExpectBegin()
	expectedDeleteExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), &domain.User{ID: userID})

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryUpdate_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{ID: int32(11), Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlUsersRepository(db)

	expectedUpdateExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`passwordHash`=?,`isAdmin`=? WHERE `id` = ?")).
			WithArgs("userName", "hash", false, 11)
	}

	mock.ExpectBegin()
	expectedUpdateExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.Update(context.Background(), &user)

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlUsersRepositoryUpdate_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{ID: int32(11), Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlUsersRepository(db)

	expectedUpdateExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`passwordHash`=?,`isAdmin`=? WHERE `id` = ?")).
			WithArgs("userName", "hash", false, 11)
	}

	mock.ExpectBegin()
	expectedUpdateExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `id` = ? LIMIT 1")).
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows(userColumns).AddRow(11, "user", "", false))

	err := repo.Update(context.Background(), &user)

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}
