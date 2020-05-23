package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestListsService(t *testing.T) {
	listColumns := []string{"id", "name", "userId"}
	listItemsColumns := []string{"id", "listId", "title", "description"}

	userId := int32(11)

	l := models.List{
		Name: "list",
	}

	i := models.ListItem{
		ListID:      11,
		Title:       "title",
		Description: "description",
	}

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	defer db.Close()

	svc := NewDefaultListsService(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`) VALUES (?,?)")).
			WithArgs(l.Name, userId)
	}

	t.Run("AddUserList() should return an error if insert fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserList(userId, &l)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserList() should insert the new list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserList(userId, &l)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, 22)
	}

	t.Run("RemoveUserList() should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserList(11, 22)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserList() should delete the user list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedRemoveListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserList(11, 22)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name` = ?, `userId` = ? WHERE `lists`.`id` = ?")).
			WithArgs(l.Name, userId, 11)
	}

	t.Run("UpdateUserList() should return an error if delete fails", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.UpdateUserList(11, userId, &l)

		appErrors.CheckUnexpectedError(t, err, "Error updating list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserList() should update the list", func(t *testing.T) {
		mock.ExpectBegin()
		expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists`  WHERE `lists`.`id` = ? ORDER BY `lists`.`id` ASC LIMIT 1")).
			WithArgs(11).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, l.Name, userId))

		err := svc.UpdateUserList(11, userId, &l)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedGetListQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?)")).
			WithArgs(11, userId)
	}

	t.Run("GetSingleUserList() should return an error if the query fails", func(t *testing.T) {
		dto := dtos.GetSingleListResultDto{}

		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetSingleUserList(11, userId, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetSingleUserList() should get a single list", func(t *testing.T) {
		dto := dtos.GetSingleListResultDto{}

		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(int32(11)).
			WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(22, 11, "title", "description"))

		err := svc.GetSingleUserList(11, userId, &dto)

		assert.Equal(t, "list", dto.Name)
		assert.Equal(t, len(dto.ListItems), 1)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name FROM `lists` WHERE (`lists`.`userId` = ?)")).
			WithArgs(userId)
	}

	t.Run("GetUserLists() should return an error if the query fails", func(t *testing.T) {
		dto := []dtos.GetListsResultDto{}

		expectedGetListsQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUserLists(userId, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user lists", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetUserLists() should return the user lists", func(t *testing.T) {
		dto := []dtos.GetListsResultDto{}

		expectedGetListsQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list1", userId).AddRow(12, "list2", userId))

		err := svc.GetUserLists(userId, &dto)

		assert.Equal(t, len(dto), 2)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT `listItems`.* FROM `listItems` JOIN lists on listItems.listId=lists.id WHERE (`lists`.`id` = ?) AND (`lists`.`userId` = ?) AND (`listItems`.`id` = ?)")).
			WithArgs(15, userId, 5)
	}

	t.Run("GetUserListItem() should return an error if the query fails", func(t *testing.T) {

		dto := dtos.GetItemResultDto{}

		expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.GetUserListItem(5, 15, userId, &dto)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("GetUserListItem() should get a single item", func(t *testing.T) {
		dto := dtos.GetItemResultDto{}

		expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(22, 11, "title", "description"))

		err := svc.GetUserListItem(5, 15, userId, &dto)

		assert.Equal(t, "title", dto.Title)
		assert.Equal(t, "description", dto.Description)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		_, err := svc.AddUserListItem(userId, &i)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	expectedListItemOpPreviousQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems`  WHERE (`listId` IN (?))")).
			WithArgs(11)
	}
	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`title`,`description`) VALUES (?,?,?)")).
			WithArgs(i.ListID, i.Title, i.Description)
	}

	t.Run("AddUserListItem() should return an error if insert fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		_, err := svc.AddUserListItem(userId, &i)

		appErrors.CheckUnexpectedError(t, err, "Error inserting list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("AddUserListItem() should insert the new list item", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))
		mock.ExpectBegin()
		expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
		mock.ExpectCommit()

		id, err := svc.AddUserListItem(userId, &i)

		assert.Equal(t, int32(12), id)
		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.RemoveUserListItem(i.ID, l.ID, userId)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE (`listItems`.`id` = ?) AND (`listItems`.`listId` = ?)")).
			WithArgs(i.ID, i.ListID)
	}

	t.Run("RemoveUserListItem() should return an error if delete fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.RemoveUserListItem(i.ID, l.ID, userId)

		appErrors.CheckUnexpectedError(t, err, "Error deleting user list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("RemoveUserListItem() should delete the user list item", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))
		mock.ExpectBegin()
		expectedRemoveListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := svc.RemoveUserListItem(i.ID, l.ID, userId)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserListItem() should return an error if getting the list fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnError(fmt.Errorf("some error"))

		err := svc.UpdateUserListItem(i.ID, l.ID, userId, &i)

		appErrors.CheckUnexpectedError(t, err, "Error getting user list", "some error")

		checkMockExpectations(t, mock)
	})

	expectedUpdateListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `listItems` SET `listId` = ?, `title` = ?, `description` = ? WHERE `listItems`.`id` = ?")).
			WithArgs(i.ListID, i.Title, i.Description, i.ID)
	}

	t.Run("UpdateUserListItem() should return an error if delete fails", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))
		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		err := svc.UpdateUserListItem(i.ID, l.ID, userId, &i)

		appErrors.CheckUnexpectedError(t, err, "Error updating list item", "some error")

		checkMockExpectations(t, mock)
	})

	t.Run("UpdateUserListItem() should update the list", func(t *testing.T) {
		expectedGetListQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		expectedListItemOpPreviousQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(11, "list", userId))

		mock.ExpectBegin()
		expectedUpdateListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`id` = ? ORDER BY `listItems`.`id` ASC LIMIT 1")).
			WithArgs(12).
			WillReturnRows(sqlmock.NewRows(listColumns).AddRow(12, l.Name, userId))

		err := svc.UpdateUserListItem(i.ID, l.ID, userId, &i)

		assert.Nil(t, err)

		checkMockExpectations(t, mock)
	})
}
