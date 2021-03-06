//+build !e2e

package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	listColumns      = []string{"id", "name", "userId", "itemsCount"}
	listItemsColumns = []string{"id", "listId", "title", "description"}
)

func TestMySqlListsRepositoryFindListByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("should not return a list if it does not exist", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(listColumns))

		res, err := repo.FindListByID(listID, userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindListByID(listID, userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user if it exists", func(t *testing.T) {
		expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list1", userID, int32(3)))

		res, err := repo.FindListByID(listID, userID)

		require.NotNil(t, res)
		assert.Equal(t, listID, res.ID)
		assert.Equal(t, domain.ListName("list1"), res.Name)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, int32(3), res.ItemsCount)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryGetAllLists(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	userID := int32(1)

	repo := NewMySqlListsRepository(db)

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetListsQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAllLists(userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should return the user lists", func(t *testing.T) {
		expectedGetListsQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(int32(11), "list1", userID, int32(3)).AddRow(int32(12), "list2", userID, int32(4)))

		res, err := repo.GetAllLists(userID)

		assert.Nil(t, err)
		require.NotNil(t, res)
		require.Equal(t, 2, len(res))
		assert.Equal(t, int32(11), res[0].ID)
		assert.Equal(t, domain.ListName("list1"), res[0].Name)
		assert.Equal(t, userID, res[0].UserID)
		assert.Equal(t, int32(3), res[0].ItemsCount)
		assert.Equal(t, int32(12), res[1].ID)
		assert.Equal(t, domain.ListName("list2"), res[1].Name)
		assert.Equal(t, userID, res[1].UserID)
		assert.Equal(t, int32(4), res[1].ItemsCount)

		checkMockExpectations(t, mock)
	})

}

func TestMySqlListsRepositoryCreateList(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	list := domain.List{UserID: 1, Name: "list1"}

	repo := NewMySqlListsRepository(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`,`itemsCount`) VALUES (?,?,?)")).
			WithArgs(list.Name, list.UserID, list.ItemsCount)
	}

	t.Run("should return an error if create fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.CreateList(&list)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should create the new list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		err := repo.CreateList(&list)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryDeleteList(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	userID := int32(1)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DeleteList(listID, userID)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should delete the user list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteList(listID, userID)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryUpdate(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	list := domain.List{ID: 11, UserID: 1, Name: "list1"}

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ? WHERE `lists`.`id` = ?")).
			WithArgs("list1", int32(11))
	}

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.UpdateList(&list)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should update the list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.UpdateList(&list)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryIncrementListCounter(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount` = itemsCount + ? WHERE (`lists`.`id` = ?)")).
			WithArgs(1, int32(11))
	}

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.IncrementListCounter(int32(11))

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should increment the items counter", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.IncrementListCounter(int32(11))

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryDecrementListCounter(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount` = itemsCount - ? WHERE (`lists`.`id` = ?)")).
			WithArgs(1, int32(11))
	}

	t.Run("should return an error if the update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DecrementListCounter(int32(11))

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should increment the items counter", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DecrementListCounter(int32(11))

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryFindListItemByID(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE (`listItems`.`id` = ?) AND (`listItems`.`listId` = ?) AND (`listItems`.`userId` = ?) ORDER BY `listItems`.`id` ASC LIMIT 1")).
			WithArgs(itemID, listID, userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.FindListItemByID(itemID, listID, userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should not return an item if it does not exist", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns))

		res, err := repo.FindListItemByID(itemID, listID, userID)

		assert.Nil(t, res)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("should get an item", func(t *testing.T) {
		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description"))

		res, err := repo.FindListItemByID(itemID, listID, userID)

		require.NotNil(t, res)
		assert.Equal(t, domain.ItemTitle("title"), res.Title)
		assert.Equal(t, "description", res.Description)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryGetAllItems(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)

	expectedGetAllItemsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE (`listItems`.`listId` = ?) AND (`listItems`.`userId` = ?)")).
			WithArgs(listID, userID)
	}

	t.Run("should return an error if the query fails", func(t *testing.T) {
		expectedGetAllItemsQuery().WillReturnError(fmt.Errorf("some error"))

		res, err := repo.GetAllListItems(listID, userID)

		assert.Nil(t, res)
		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should get all the items", func(t *testing.T) {
		expectedGetAllItemsQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(int32(111), listID, "title1", "desc1").AddRow(int32(112), listID, "title2", "desc2"))

		res, err := repo.GetAllListItems(listID, userID)

		require.NotNil(t, res)
		require.Equal(t, 2, len(res))
		assert.Equal(t, domain.ItemTitle("title1"), res[0].Title)
		assert.Equal(t, "desc1", res[0].Description)
		assert.Equal(t, domain.ItemTitle("title2"), res[1].Title)
		assert.Equal(t, "desc2", res[1].Description)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryCreateListItem(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	item := domain.ListItem{ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`) VALUES (?,?,?,?)")).
			WithArgs(int32(11), int32(1), "title", "desc")
	}

	t.Run("should return an error if create fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.CreateListItem(&item)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should create the new list item", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		err := repo.CreateListItem(&item)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryDeleteListItem(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE (`listItems`.`id` = ?) AND (`listItems`.`listId` = ?) AND (`listItems`.`userId` = ?)")).
			WithArgs(itemID, listID, userID)
	}

	t.Run("should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.DeleteListItem(itemID, listID, userID)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should delete the item", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteListItem(itemID, listID, userID)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func TestMySqlListsRepositoryUpdateListItem(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	repo := NewMySqlListsRepository(db)

	item := domain.ListItem{ID: 111, ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedUpdateListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `listItems` SET `listId` = ?, `userId` = ?, `title` = ?, `description` = ? WHERE `listItems`.`id` = ?")).
			WithArgs(int32(11), int32(1), "title", "desc", int32(111))
	}

	t.Run("should return an error if update fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := repo.UpdateListItem(&item)

		assert.EqualError(t, err, "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("should update the list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`id` = ? ORDER BY `listItems`.`id` ASC LIMIT 1")).
			WithArgs(int32(111)).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(int32(111), int32(11), "title", "desc"))

		err := repo.UpdateListItem(&item)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
