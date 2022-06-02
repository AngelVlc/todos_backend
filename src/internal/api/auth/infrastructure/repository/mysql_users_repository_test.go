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
