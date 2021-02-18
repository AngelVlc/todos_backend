package repositories

import (
	"fmt"
	"regexp"
	"testing"

	appErrors "github.com/AngelVlc/todos/errors"
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

func TestUsersServiceFindByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	svc := NewDefaultUsersRepository(db)

	expectedFindByIdQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (`users`.`id` = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
			WithArgs(int32(11))
	}

	t.Run("should not return a user if it does not exist", func(t *testing.T) {
		expectedFindByIdQuery().WillReturnRows(sqlmock.NewRows(columns))

		dto, err := svc.FindByID(11)

		assert.Nil(t, dto)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIdQuery().WillReturnError(fmt.Errorf("some error"))

		dto, err := svc.FindByID(11)

		assert.Nil(t, dto)
		appErrors.CheckUnexpectedError(t, err, "Error getting user by user id", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIdQuery().WillReturnRows(sqlmock.NewRows(columns).AddRow(5, user, "", true))

		dto, err := svc.FindByID(11)

		require.NotNil(t, dto)
		assert.Equal(t, user, dto.Name)
		assert.Equal(t, true, dto.IsAdmin)
		assert.Equal(t, int32(5), dto.ID)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
