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
)

func TestListsRepositoryInsert(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	list := models.List{UserID: 1, Name: "list1"}

	repo := NewDefaultListsRepository(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`) VALUES (?,?)")).
			WithArgs(list.Name, list.UserID)
	}

	t.Run("should return an error if insert fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := repo.Insert(&list)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should insert the new list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := repo.Insert(&list)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}
