//+build !e2e

package infrastructure

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/internal/api/shared/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	counterColumns = []string{"id", "name", "value"}
)

func TestMySqlCountersRepositoryFindByName(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	repo := NewMySqlCountersRepository(db)

	name := "counter-name"

	expectedFindByNameQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE `counters`.`name` = ? ORDER BY `counters`.`id` LIMIT 1")).
			WithArgs("counter-name")
	}

	t.Run("should not return the counter if it does not exist", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(counterColumns))

		c, err := repo.FindByName(name)

		assert.Nil(t, c)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnError(fmt.Errorf("some error"))

		c, err := repo.FindByName(name)

		assert.Nil(t, c)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the counter if it exists", func(t *testing.T) {
		expectedFindByNameQuery().WillReturnRows(sqlmock.NewRows(counterColumns).AddRow(int32(1), "counter-name", 10))

		c, err := repo.FindByName(name)

		assert.NotNil(t, c)
		assert.Equal(t, int32(1), c.ID)
		assert.Equal(t, "counter-name", c.Name)
		assert.Equal(t, int32(10), c.Value)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlCountersRepositoryCreate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	counter := domain.Counter{Name: "counter", Value: 0}

	repo := NewMySqlCountersRepository(db)

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `counters` (`name`,`value`) VALUES (?,?)")).
			WithArgs(counter.Name, counter.Value)
	}

	t.Run("should return an error if creating the new counter fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Create(&counter)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should create the new counter", func(t *testing.T) {
		result := sqlmock.NewResult(12, 1)

		mock.ExpectBegin()
		expectedInsertExec().WillReturnResult(result)
		mock.ExpectCommit()

		err := repo.Create(&counter)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlCountersRepositoryUpdate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	counter := domain.Counter{ID: int32(11), Name: "counter", Value: 100}

	repo := NewMySqlCountersRepository(db)

	expectedUpdateExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `counters` SET `name`=?,`value`=? WHERE `id` = ?")).
			WithArgs("counter", int32(100), 11)
	}

	t.Run("should return an error if update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Update(&counter)

		assert.EqualError(t, err, "some error")
		checkMockExpectations(t, mock)
	})

	t.Run("should update the counter if the update doesn't fail", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE `id` = ? ORDER BY `counters`.`id` LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(counterColumns).AddRow(11, "counter", int32(100)))

		err := repo.Update(&counter)

		assert.Nil(t, err)
		checkMockExpectations(t, mock)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
