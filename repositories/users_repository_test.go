package repositories

import (
	"fmt"
	"regexp"
	"testing"

	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	columns     = []string{"id", "name", "password_hash", "is_admin"}
	hasshedPass = "hassedPassword"
	user        = "user"
)

func TestUsersRepositoryFindByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultUsersRepository(db)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11))
	}

	t.Run("should not return a user if it does not exist", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(columns))

		dto, err := repo.FindByID(11)

		assert.Nil(t, dto)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		dto, err := repo.FindByID(11)

		assert.Nil(t, dto)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user id", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		dto, err := repo.FindByID(11)

		require.NotNil(t, dto)
		assert.Equal(t, user, dto.Name)
		assert.Equal(t, true, dto.IsAdmin)
		assert.Equal(t, int32(5), dto.ID)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestUsersRepositoryFindByName(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultUsersRepository(db)

	expectedFindByNameQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`name` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(user)
	}

	t.Run("should not return a user if it does not exist", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(columns))

		u, err := repo.FindByName(user)

		assert.Nil(t, u)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnError(fmt.Errorf("some error"))

		u, err := repo.FindByName(user)

		assert.Nil(t, u)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user name", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		u, err := repo.FindByName(user)

		assert.NotNil(t, u)
		assert.Equal(t, u.ID, int32(5))
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestUsersRepositoryInsert(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	user := models.User{Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewDefaultUsersRepository(db)

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`name`,`password_hash`,`is_admin`) VALUES (?,?,?)")).
			WithArgs(user.Name, user.PasswordHash, user.IsAdmin)
	}

	t.Run("should return an error if inserting the new user fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := repo.Insert(&user)

		appErrors.CheckUnexpectedError(t, err, "Error inserting in the database", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should insert the new user", func(t *testing.T) {
		var affected int64
		result := sqlmock.NewResult(12, affected)

		mock.ExpectBegin()
		expectedInsertExec().WillReturnResult(result)
		mock.ExpectCommit()

		id, err := repo.Insert(&user)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestUsersRepositoryRemove(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultUsersRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE (`users`.`id` = ?)")).
			WithArgs(11)
	}

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Remove(11)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user", "some error")
		checkMockExpectations(t, mock)
	})

	t.Run("should delete the user", func(t *testing.T) {
		mock.ExpectBegin()
		expectedDeleteExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Remove(11)

		assert.Nil(t, err)
		checkMockExpectations(t, mock)
	})
}

func TestUsersRepositoryUpdate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	user := models.User{ID: int32(11), Name: "userName", PasswordHash: "hash", IsAdmin: false}

	repo := NewDefaultUsersRepository(db)

	expectedUpdateExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name` = ?, `password_hash` = ?, `is_admin` = ? WHERE `users`.`id` = ?")).
			WithArgs("userName", "hash", false, 11)
	}

	t.Run("should return an error if update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Update(&user)

		assert.NotNil(t, err)
		checkMockExpectations(t, mock)
	})

	t.Run("should update the user if the update doesn't fail", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`  WHERE `users`.`id` = ? ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(11, "user", "", false))

		err := repo.Update(&user)

		assert.Nil(t, err)
		checkMockExpectations(t, mock)
	})

}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
