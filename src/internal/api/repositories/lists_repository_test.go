//+build !e2e

package repositories

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/internal/api/models"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/infrastructure/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	listColumns      = []string{"id", "name", "userId"}
	listItemsColumns = []string{"id", "listId", "title", "description"}
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

	t.Run("should return an error if create fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		id, err := repo.Create(&list)

		assert.Equal(t, int32(-1), id)
		appErrors.CheckUnexpectedError(t, err, "Error creating list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should create the new list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := repo.Create(&list)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestListsRepositoryRemove(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	userID := int32(1)
	listID := int32(11)

	repo := NewDefaultListsRepository(db)

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Delete(listID, userID)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should delete the user list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(listID, userID)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestListsRepositoryUpdate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultListsRepository(db)

	list := models.List{ID: 11, UserID: 1, Name: "list1"}

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ?, `userId` = ? WHERE `lists`.`id` = ?")).
			WithArgs("list1", int32(1), int32(11))
	}

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.Update(&list)

		appErrors.CheckUnexpectedError(t, err, "Error updating list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should update the list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`  WHERE `lists`.`id` = ? ORDER BY `lists`.`id` ASC LIMIT 1")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(int32(11), "list1", int32(1)))

		err := repo.Update(&list)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestListsRepositoryFindByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewDefaultListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("should not return a list if it does not exist", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(columns))

		res, err := repo.FindByID(listID, userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindByID(listID, userID)

		assert.Nil(t, res)
		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list1", userID))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(listID).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(int32(111), listID, "title", "description"))

		res, err := repo.FindByID(listID, userID)

		require.NotNil(t, res)
		assert.Equal(t, listID, res.ID)
		assert.Equal(t, "list1", res.Name)
		assert.Equal(t, userID, res.UserID)
		require.Equal(t, 1, len(res.ListItems))
		assert.Equal(t, "title", res.ListItems[0].Title)
		assert.Equal(t, "description", res.ListItems[0].Description)
		assert.Equal(t, listID, res.ListItems[0].ListID)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestListsRepositoryGetAll(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	userID := int32(1)

	repo := NewDefaultListsRepository(db)

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetListsQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAll(userID)

		assert.Nil(t, res)
		appErrors.CheckUnexpectedError(t, err, "Error getting user lists", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user lists", func(t *testing.T) {
		expectedGetListsQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(int32(11), "list1", userID).AddRow(int32(12), "list2", userID))

		res, err := repo.GetAll(userID)

		assert.Nil(t, err)
		require.NotNil(t, res)
		require.Equal(t, 2, len(res))
		assert.Equal(t, int32(11), res[0].ID)
		assert.Equal(t, "list1", res[0].Name)
		assert.Equal(t, userID, res[0].UserID)
		assert.Equal(t, int32(12), res[1].ID)
		assert.Equal(t, "list2", res[1].Name)
		assert.Equal(t, userID, res[1].UserID)

		checkMockExpectations(t, mock)
	})

}
