//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	listColumns      = []string{"id", "name", "userId", "itemsCount"}
	listItemsColumns = []string{"id", "listId", "title", "description", "position"}
)

func TestMySqlListsRepository_FindList_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	expectedFindByIDQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.FindList(context.Background(), &domain.ListEntity{ID: listID, UserID: userID})

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_FindList_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listID := int32(11)
	userID := int32(1)

	expectedFindByIDQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	listName, _ := domain.NewListNameValueObject("list1")

	expectedFindByIDQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(listID, "list1", userID, int32(3)))

	res, err := repo.FindList(context.Background(), &domain.ListEntity{ID: listID, UserID: userID})

	require.NotNil(t, res)
	assert.Equal(t, listID, res.ID)
	assert.Equal(t, listName, res.Name)
	assert.Equal(t, userID, res.UserID)
	assert.Equal(t, int32(3), res.ItemsCount)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_ExistsList_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listName, _ := domain.NewListNameValueObject("list name")
	userID := int32(1)
	list := &domain.ListEntity{Name: listName, UserID: userID}

	expectedExistsListQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `lists` WHERE `lists`.`name` = ? AND `lists`.`userId` = ?")).
			WithArgs(listName, userID)
	}

	expectedExistsListQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.ExistsList(context.Background(), list)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_ExistsList_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listName, _ := domain.NewListNameValueObject("list name")
	userID := int32(1)
	list := &domain.ListEntity{Name: listName, UserID: userID}

	expectedExistsListQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `lists` WHERE `lists`.`name` = ? AND `lists`.`userId` = ?")).
			WithArgs(listName, userID)
	}

	expectedExistsListQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	res, err := repo.ExistsList(context.Background(), list)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllLists_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)

	repo := NewMySqlListsRepository(db)

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`userId` = ?")).
			WithArgs(userID)
	}

	expectedGetListsQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllLists(context.Background(), userID)

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllLists_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)

	repo := NewMySqlListsRepository(db)

	expectedGetListsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `lists` WHERE `lists`.`userId` = ?")).
			WithArgs(userID)
	}

	expectedGetListsQuery().WillReturnRows(sqlmock.NewRows(listColumns).AddRow(int32(11), "list1", userID, int32(3)).AddRow(int32(12), "list2", userID, int32(4)))

	res, err := repo.GetAllLists(context.Background(), userID)

	list1Name, _ := domain.NewListNameValueObject("list1")
	list2Name, _ := domain.NewListNameValueObject("list2")

	assert.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, list1Name, res[0].Name)
	assert.Equal(t, userID, res[0].UserID)
	assert.Equal(t, int32(3), res[0].ItemsCount)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, list2Name, res[1].Name)
	assert.Equal(t, userID, res[1].UserID)
	assert.Equal(t, int32(4), res[1].ItemsCount)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateList_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListEntity{UserID: 1, Name: listName}

	repo := NewMySqlListsRepository(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`,`itemsCount`) VALUES (?,?,?)")).
			WithArgs(list.Name, list.UserID, list.ItemsCount)
	}

	mock.ExpectBegin()
	expectedInsertListExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.CreateList(context.Background(), &list)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateList_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListEntity{UserID: 1, Name: listName}

	repo := NewMySqlListsRepository(db)

	expectedInsertListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`name`,`userId`,`itemsCount`) VALUES (?,?,?)")).
			WithArgs(list.Name, list.UserID, list.ItemsCount)
	}

	mock.ExpectBegin()
	expectedInsertListExec().WillReturnResult(sqlmock.NewResult(12, 0))
	mock.ExpectCommit()

	err := repo.CreateList(context.Background(), &list)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_Deleting_The_ListItems_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	expectedRemoveListItemsExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	mock.ExpectBegin()
	expectedRemoveListItemsExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteList(context.Background(), listID, userID)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_Deleting_The_List_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	expectedRemoveListItemsExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	mock.ExpectBegin()
	expectedRemoveListItemsExec().WillReturnResult(sqlmock.NewResult(0, 0))
	expectedRemoveListExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteList(context.Background(), listID, userID)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteList_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	userID := int32(1)
	listID := int32(11)

	repo := NewMySqlListsRepository(db)

	expectedRemoveListItemsExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	expectedRemoveListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `lists` WHERE `lists`.`id` = ? AND `lists`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	mock.ExpectBegin()
	expectedRemoveListItemsExec().WillReturnResult(sqlmock.NewResult(0, 0))
	expectedRemoveListExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteList(context.Background(), listID, userID)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_Update_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListEntity{ID: 11, UserID: 1, Name: listName}

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name`=? WHERE `id` = ?")).
			WithArgs("list1", int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.UpdateList(context.Background(), &list)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_Update_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	listName, _ := domain.NewListNameValueObject("list1")
	list := domain.ListEntity{ID: 11, UserID: 1, Name: listName}

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `name`=? WHERE `id` = ?")).
			WithArgs("list1", int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateList(context.Background(), &list)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_IncrementListCounter_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=itemsCount + ? WHERE `lists`.`id` = ?")).
			WithArgs(1, int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.IncrementListCounter(context.Background(), int32(11))

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_IncrementListCounter_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=itemsCount + ? WHERE `lists`.`id` = ?")).
			WithArgs(1, int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.IncrementListCounter(context.Background(), int32(11))

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DecrementListCounter_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=itemsCount - ? WHERE `lists`.`id` = ?")).
			WithArgs(1, int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DecrementListCounter(context.Background(), int32(11))

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DecrementListCounter_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	expectedUpdateListExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `itemsCount`=itemsCount - ? WHERE `lists`.`id` = ?")).
			WithArgs(1, int32(11))
	}

	mock.ExpectBegin()
	expectedUpdateListExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DecrementListCounter(context.Background(), int32(11))

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_FindListItem_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	listItem := domain.ListItemEntity{ID: itemID, ListID: listID, UserID: userID}

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`id` = ? AND `listItems`.`listId` = ? AND `listItems`.`userId` = ? LIMIT 1")).
			WithArgs(itemID, listID, userID)
	}

	expectedGetItemQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.FindListItem(context.Background(), &listItem)

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_FindListItem_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	listItem := domain.ListItemEntity{ID: itemID, ListID: listID, UserID: userID}

	expectedGetItemQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`id` = ? AND `listItems`.`listId` = ? AND `listItems`.`userId` = ? LIMIT 1")).
			WithArgs(itemID, listID, userID)
	}

	expectedGetItemQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(itemID, listID, "title", "description", 0))

	res, err := repo.FindListItem(context.Background(), &listItem)

	require.NotNil(t, res)
	assert.Equal(t, domain.ItemTitleValueObject("title"), res.Title)
	assert.Equal(t, domain.ItemDescriptionValueObject("description"), res.Description)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllItems_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)

	expectedGetAllItemsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ? ORDER BY position")).
			WithArgs(listID, userID)
	}

	expectedGetAllItemsQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllListItems(context.Background(), listID, userID)

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetAllItems_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)

	expectedGetAllItemsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ? ORDER BY position")).
			WithArgs(listID, userID)
	}

	expectedGetAllItemsQuery().WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(int32(111), listID, "title1", "desc1", 0).AddRow(int32(112), listID, "title2", "desc2", 1))

	res, err := repo.GetAllListItems(context.Background(), listID, userID)

	require.NotNil(t, res)
	require.Equal(t, 2, len(res))
	assert.Equal(t, domain.ItemTitleValueObject("title1"), res[0].Title)
	assert.Equal(t, domain.ItemDescriptionValueObject("desc1"), res[0].Description)
	assert.Equal(t, int32(0), res[0].Position)
	assert.Equal(t, domain.ItemTitleValueObject("title2"), res[1].Title)
	assert.Equal(t, domain.ItemDescriptionValueObject("desc2"), res[1].Description)
	assert.Equal(t, int32(1), res[1].Position)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateListItem_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item := domain.ListItemEntity{ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`,`position`) VALUES (?,?,?,?,?)")).
			WithArgs(int32(11), int32(1), "title", "desc", int32(0))
	}

	mock.ExpectBegin()
	expectedInsertListItemExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.CreateListItem(context.Background(), &item)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_CreateListItem_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item := domain.ListItemEntity{ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedInsertListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`,`position`) VALUES (?,?,?,?,?)")).
			WithArgs(int32(11), int32(1), "title", "desc", int32(0))
	}

	mock.ExpectBegin()
	expectedInsertListItemExec().WillReturnResult(sqlmock.NewResult(12, 0))
	mock.ExpectCommit()

	err := repo.CreateListItem(context.Background(), &item)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteListItem_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`id` = ? AND `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(itemID, listID, userID)
	}

	mock.ExpectBegin()
	expectedRemoveListItemExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteListItem(context.Background(), itemID, listID, userID)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_DeleteListItem_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)
	itemID := int32(111)

	expectedRemoveListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `listItems` WHERE `listItems`.`id` = ? AND `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(itemID, listID, userID)
	}

	mock.ExpectBegin()
	expectedRemoveListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteListItem(context.Background(), itemID, listID, userID)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateListItem_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item := domain.ListItemEntity{ID: 111, ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedUpdateListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `listItems` SET `listId`=?,`userId`=?,`title`=?,`description`=?,`position`=? WHERE `id` = ?")).
			WithArgs(int32(11), int32(1), "title", "desc", int32(0), int32(111))
	}

	mock.ExpectBegin()
	expectedUpdateListItemExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.UpdateListItem(context.Background(), &item)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_UpdateListItem(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item := domain.ListItemEntity{ID: 111, ListID: 11, UserID: 1, Title: "title", Description: "desc"}

	expectedUpdateListItemExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("UPDATE `listItems` SET `listId`=?,`userId`=?,`title`=?,`description`=?,`position`=? WHERE `id` = ?")).
			WithArgs(int32(11), int32(1), "title", "desc", int32(0), int32(111))
	}

	mock.ExpectBegin()
	expectedUpdateListItemExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `listItems` WHERE `id` = ? LIMIT 1")).
		WithArgs(int32(111)).
		WillReturnRows(sqlmock.NewRows(listItemsColumns).AddRow(int32(111), int32(11), "title", "desc", 0))

	err := repo.UpdateListItem(context.Background(), &item)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_BulkUpdateListItems_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item1 := domain.ListItemEntity{ID: 1, Position: 0}
	item2 := domain.ListItemEntity{ID: 2, Position: 1}
	item3 := domain.ListItemEntity{ID: 3, Position: 2}
	items := []domain.ListItemEntity{item1, item2, item3}

	expectedUpdateListItemsExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`,`position`,`id`) VALUES (?,?,?,?,?,?),(?,?,?,?,?,?),(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `position`=VALUES(`position`)")).
			WithArgs(int32(0), int32(0), "", "", int32(0), int32(1), int32(0), int32(0), "", "", int32(1), int32(2), int32(0), int32(0), "", "", int32(2), int32(3))
	}

	mock.ExpectBegin()
	expectedUpdateListItemsExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.BulkUpdateListItems(context.Background(), items)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_BulkUpdateListItems_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	item1 := domain.ListItemEntity{ID: 1, Position: 0}
	item2 := domain.ListItemEntity{ID: 2, Position: 1}
	item3 := domain.ListItemEntity{ID: 3, Position: 2}
	items := []domain.ListItemEntity{item1, item2, item3}

	expectedUpdateListItemsExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `listItems` (`listId`,`userId`,`title`,`description`,`position`,`id`) VALUES (?,?,?,?,?,?),(?,?,?,?,?,?),(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `position`=VALUES(`position`)")).
			WithArgs(int32(0), int32(0), "", "", int32(0), int32(1), int32(0), int32(0), "", "", int32(1), int32(2), int32(0), int32(0), "", "", int32(2), int32(3))
	}

	mock.ExpectBegin()
	expectedUpdateListItemsExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.BulkUpdateListItems(context.Background(), items)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetListItemsMaxPosition_When_It_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)

	expectedGetAllItemsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT MAX(position) FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	expectedGetAllItemsQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetListItemsMaxPosition(context.Background(), listID, userID)

	assert.Equal(t, int32(-1), res)
	assert.EqualError(t, err, "some error; some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlListsRepository_GetListItemsMaxPosition_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlListsRepository(db)

	userID := int32(1)
	listID := int32(11)

	expectedGetAllItemsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT MAX(position) FROM `listItems` WHERE `listItems`.`listId` = ? AND `listItems`.`userId` = ?")).
			WithArgs(listID, userID)
	}

	expectedGetAllItemsQuery().WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(int32(3)))

	res, err := repo.GetListItemsMaxPosition(context.Background(), listID, userID)

	require.NotNil(t, res)
	require.Equal(t, int32(3), res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}
