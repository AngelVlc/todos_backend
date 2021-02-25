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

func TestListItemsRepositoryFindByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultListItemsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT `listItems`.* FROM `listItems` JOIN lists on listItems.listId=lists.id WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?) AND (`listItems`.`id` = ?)")).
			WithArgs(listID, userID, itemID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindByID(itemID, listID, userID)

		assert.Nil(t, res)
		appErrors.CheckUnexpectedError(t, err, "Error getting user list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should not return an item if it does not exist", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(columns))

		res, err := repo.FindByID(itemID, listID, userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should get an item", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		res, err := repo.FindByID(itemID, listID, userID)

		require.NotNil(t, res)
		assert.Equal(t, "title", res.Title)
		assert.Equal(t, "description", res.Description)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

}

func TestListItemsRepositoryInsert(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultListItemsRepository(db)

	item := models.ListItem{ListID: 11, Title: "title", Description: "desc"}

	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`title`,`description`) VALUES (?,?,?)")).
			WithArgs(int32(11), "title", "desc")
	}

	t.Run("should return an error if insert fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		id, err := repo.Insert(&item)

		assert.Equal(t, int32(-1), id)
		appErrors.CheckUnexpectedError(t, err, "Error inserting list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should insert the new list item", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := repo.Insert(&item)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

}
