package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/internal/api/auth/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userColumns = []string{"id", "name", "password_hash", "is_admin"}
)

func TestMySqlAuthRepositoryFindUserByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlAuthRepository(db)

	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(userID)
	}

	t.Run("should not return a user if it does not exist", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(userColumns))

		res, err := repo.FindUserByID(&userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindUserByID(&userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(userID, "userName", "hash", true))

		res, err := repo.FindUserByID(&userID)

		require.NotNil(t, res)
		assert.Equal(t, domain.AuthUserName("userName"), res.Name)
		assert.True(t, res.IsAdmin)
		assert.Equal(t, userID, res.ID)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlAuthRepositoryFindUserByName(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlAuthRepository(db)

	userName := domain.AuthUserName("userName")

	expectedFindByNameQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs("userName")
	}

	t.Run("should not return a user if it does not exist", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(userColumns))

		u, err := repo.FindUserByName(&userName)

		assert.Nil(t, u)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnError(fmt.Errorf("some error"))

		u, err := repo.FindUserByName(&userName)

		assert.Nil(t, u)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(int32(1), "userName", "hash", true))

		u, err := repo.FindUserByName(&userName)

		assert.NotNil(t, u)
		assert.Equal(t, int32(1), u.ID)
		assert.True(t, u.IsAdmin)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlAuthRepositoryGetAllUsers(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlAuthRepository(db)

	expectedGetUsersQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,is_admin FROM `users`"))
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetUsersQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAllUsers()

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the users", func(t *testing.T) {
		expectedGetUsersQuery().WillReturnRows(sqlmock.NewRows(userColumns).AddRow(11, "user1", "pass1", true).AddRow(12, "user2", "pass2", false))

		res, err := repo.GetAllUsers()

		assert.Nil(t, err)
		require.Equal(t, 2, len(res))
		assert.Equal(t, int32(11), res[0].ID)
		assert.Equal(t, domain.AuthUserName("user1"), res[0].Name)
		assert.Equal(t, "pass1", res[0].PasswordHash)
		assert.True(t, res[0].IsAdmin)
		assert.Equal(t, int32(12), res[1].ID)
		assert.Equal(t, domain.AuthUserName("user2"), res[1].Name)
		assert.Equal(t, "pass2", res[1].PasswordHash)
		assert.False(t, res[1].IsAdmin)

		checkMockExpectations(t, mock)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
