package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCountersService(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()
	svc := NewDefaultCountersService(db)

	counterName := "counter"
	columns := []string{"id", "name", "value"}

	t.Run("CreateCounterIfNotExists() should do nothing if the counter exists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, counterName, 11))

		err := svc.CreateCounterIfNotExists(counterName)

		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("CreateCounterIfNotExists() should create the counter if not exists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns))

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `counters` (`name`,`value`) VALUES (?,?)")).
			WithArgs(counterName, 0).
			WillReturnResult(result)
		mock.ExpectCommit()

		err := svc.CreateCounterIfNotExists(counterName)

		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("CreateCounterIfNotExists() should return an error if insert fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `counters` (`name`,`value`) VALUES (?,?)")).
			WithArgs(counterName, 0).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.CreateCounterIfNotExists(counterName)

		if assert.Error(t, err) {
			assert.Equal(t, "some error", err.Error())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("IncrementCounter() should return an error if the counter does not exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns))

		_, err := svc.IncrementCounter(counterName)

		if assert.Error(t, err) {
			assert.Equal(t, "error getting 'counter' counter: record not found", err.Error())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("IncrementCounter() should return an error if updating the counter fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, counterName, 11))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `counters` SET `name` = ?, `value` = ? WHERE `counters`.`id` = ?")).
			WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.IncrementCounter(counterName)

		if assert.Error(t, err) {
			assert.Equal(t, "error saving new 'counter' counter value: some error", err.Error())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("IncrementCounter() should increment the counter if there aren't errors", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE (`counters`.`name` = ?) ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(counterName).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, counterName, 11))

		var lastInsertID, affected int64
		result := sqlmock.NewResult(lastInsertID, affected)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `counters` SET `name` = ?, `value` = ? WHERE `counters`.`id` = ?")).
			WithArgs(counterName, 12, 5).
			WillReturnResult(result)

		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `counters` WHERE `counters`.`id` = ? ORDER BY `counters`.`id` ASC LIMIT 1")).
			WithArgs(5).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(5, counterName, 12))

		v, err := svc.IncrementCounter(counterName)

		assert.Equal(t, int32(12), v)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
