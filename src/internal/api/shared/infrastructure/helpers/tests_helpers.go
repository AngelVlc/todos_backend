package helpers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMockedDb(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	return mock, db
}

func CheckSqlMockExpectations(mock sqlmock.Sqlmock, t *testing.T) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
