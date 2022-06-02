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

func TestMySqlAuthRepositoryFindUserByID(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlUsersRepository(db)

	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? LIMIT 1")).
			WithArgs(userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindUser(context.Background(), &domain.User{ID: userID})

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(userID, "userName", "hash", true))

		res, err := repo.FindUser(context.Background(), &domain.User{ID: userID})

		require.NotNil(t, res)
		assert.Equal(t, domain.UserName("userName"), res.Name)
		assert.True(t, res.IsAdmin)
		assert.Equal(t, userID, res.ID)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}
