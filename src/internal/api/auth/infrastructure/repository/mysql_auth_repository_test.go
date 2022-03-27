//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedDomain "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userColumns         = []string{"id", "name", "passwordHash", "isAdmin"}
	refreshTokenColumns = []string{"id", "userId", "refreshToken", "expirationDate"}
)

func TestMySqlAuthRepositoryExistsUser(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	userName := domain.UserName("userName")

	expectedExistsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`name` = ?")).
			WithArgs("userName")
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedExistsQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.ExistsUser(context.Background(), userName)

		assert.False(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return true if the user exists", func(t *testing.T) {
		expectedExistsQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		res, err := repo.ExistsUser(context.Background(), userName)

		assert.True(t, res)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryFindUserByID(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? LIMIT 1")).
			WithArgs(userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindUserByID(context.Background(), userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(userID, "userName", "hash", true))

		res, err := repo.FindUserByID(context.Background(), userID)

		require.NotNil(t, res)
		assert.Equal(t, domain.UserName("userName"), res.Name)
		assert.True(t, res.IsAdmin)
		assert.Equal(t, userID, res.ID)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryFindUserByName(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	userName := domain.UserName("userName")

	expectedFindByNameQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`name` = ? LIMIT 1")).
			WithArgs("userName")
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnError(fmt.Errorf("some error"))

		u, err := repo.FindUserByName(context.Background(), userName)

		assert.Nil(t, u)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(int32(1), "userName", "hash", true))

		u, err := repo.FindUserByName(context.Background(), userName)

		assert.NotNil(t, u)
		assert.Equal(t, int32(1), u.ID)
		assert.True(t, u.IsAdmin)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryGetAllUsers(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expectedGetUsersQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,isAdmin FROM `users`"))
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetUsersQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAllUsers(context.Background())

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the users", func(t *testing.T) {
		expectedGetUsersQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(11, "user1", "pass1", true).AddRow(12, "user2", "pass2", false))

		res, err := repo.GetAllUsers(context.Background())

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
	})
}

func TestMySqlAuthRepositoryCreateUser(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlAuthRepository(db)

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`passwordHash`,`isAdmin`) VALUES (?,?,?)")).
			WithArgs(user.Name, user.PasswordHash, user.IsAdmin)
	}

	t.Run("should return an error if creating the new user fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.CreateUser(context.Background(), &user)

		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should create the new user", func(t *testing.T) {
		result := sqlmock.NewResult(12, 1)

		mock.ExpectBegin()
		expectedInsertExec().WillReturnResult(result)
		mock.ExpectCommit()

		err := repo.CreateUser(context.Background(), &user)

		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryDeleteUser(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
			WithArgs(1)
	}

	userID := int32(1)

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DeleteUser(context.Background(), userID)

		assert.EqualError(t, err, "some error")
		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should delete the user", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteUser(context.Background(), userID)

		assert.Nil(t, err)
		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryUpdateUser(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	user := domain.User{ID: int32(11), Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewMySqlAuthRepository(db)

	expectedUpdateExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`passwordHash`=?,`isAdmin`=? WHERE `id` = ?")).
			WithArgs("userName", "hash", false, 11)
	}

	t.Run("should return an error if update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.UpdateUser(context.Background(), &user)

		assert.EqualError(t, err, "some error")
		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should update the user if the update doesn't fail", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `id` = ? LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(userColumns).AddRow(11, "user", "", false))

		err := repo.UpdateUser(context.Background(), &user)

		assert.Nil(t, err)
		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryFindRefreshTokenForUser(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	rt := "rt"
	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `refresh_tokens` WHERE `refresh_tokens`.`userId` = ? AND `refresh_tokens`.`refreshToken` = ? LIMIT 1")).
			WithArgs(userID, rt)
	}

	t.Run("should not return a refresh token if it does not exist", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(refreshTokenColumns))

		res, err := repo.FindRefreshTokenForUser(context.Background(), rt, userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindRefreshTokenForUser(context.Background(), rt, userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the refresh token if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(refreshTokenColumns).AddRow(int32(111), userID, rt, time.Now()))

		res, err := repo.FindRefreshTokenForUser(context.Background(), rt, userID)

		require.NotNil(t, res)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, int32(111), res.ID)
		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryCreateRefreshTokenIfNotExist(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expDate, _ := time.Parse("2021-Jan-01", "2014-Feb-04")
	rt := domain.RefreshToken{UserID: 1, RefreshToken: "rt", ExpirationDate: expDate}

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `refresh_tokens` (`userId`,`refreshToken`,`expirationDate`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
			WithArgs(rt.UserID, rt.RefreshToken, rt.ExpirationDate)
	}

	t.Run("should return an error if creating the new refresh token fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.CreateRefreshTokenIfNotExist(context.Background(), &rt)

		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should create the new refresh token", func(t *testing.T) {
		result := sqlmock.NewResult(12, 1)

		mock.ExpectBegin()
		expectedInsertExec().WillReturnResult(result)
		mock.ExpectCommit()

		err := repo.CreateRefreshTokenIfNotExist(context.Background(), &rt)

		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthDeleteExpiredRefreshTokens(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expectedDelete := func(expTime time.Time) *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE expirationDate <= ?")).
			WithArgs(expTime)
	}

	t.Run("should return an error if creating the new refresh token fails", func(t *testing.T) {
		now := time.Now()
		mock.ExpectBegin()
		expectedDelete(now).WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DeleteExpiredRefreshTokens(context.Background(), now)

		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should create the new refresh token", func(t *testing.T) {
		now := time.Now()

		mock.ExpectBegin()
		expectedDelete(now).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteExpiredRefreshTokens(context.Background(), now)

		assert.Nil(t, err)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryGetAllRefreshTokens(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	paginationInfo := sharedDomain.NewPaginationInfo(10, 10, "expirationDate", sharedDomain.OrderAsc)

	expectedGetQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT id,userId,expirationDate FROM `refresh_tokens` ORDER BY expirationDate asc LIMIT 10 OFFSET 10")
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAllRefreshTokens(context.Background(), paginationInfo)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should return the users", func(t *testing.T) {
		columns := []string{"id", "userId", "expirationDate"}
		now := time.Now()
		expectedGetQuery().WillReturnRows(sqlmock.NewRows(columns).AddRow(11, 1, now))

		res, err := repo.GetAllRefreshTokens(context.Background(), paginationInfo)

		assert.Nil(t, err)
		require.Equal(t, 1, len(res))
		assert.Equal(t, int32(11), res[0].ID)
		assert.Equal(t, int32(1), res[0].UserID)
		assert.Equal(t, now, res[0].ExpirationDate)

		helpers.CheckSqlMockExpectations(mock, t)
	})
}

func TestMySqlAuthRepositoryDeleteRefreshTokensByID(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlAuthRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `refresh_tokens` WHERE `refresh_tokens`.`id` IN (?,?,?)")).
			WithArgs(int32(1), int32(2), int32(3))
	}

	ids := []int32{int32(1), int32(2), int32(3)}

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DeleteRefreshTokensByID(context.Background(), ids)

		assert.EqualError(t, err, "some error")
		helpers.CheckSqlMockExpectations(mock, t)
	})

	t.Run("should delete the refresh tokens", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteRefreshTokensByID(context.Background(), ids)

		assert.Nil(t, err)
		helpers.CheckSqlMockExpectations(mock, t)
	})
}
