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

var categoryColumns = []string{"id", "name", "description"}

func TestMySqlCategoriesRepository_FindCategory_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlCategoriesRepository(db)

	categoryID := int32(11)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(categoryID).
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.FindCategory(context.Background(), domain.CategoryEntity{ID: categoryID})

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_FindCategory_WhenTheQueryDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlCategoriesRepository(db)

	categoryID := int32(11)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows(categoryColumns).
			AddRow(categoryID, "name", "description"))

	res, err := repo.FindCategory(context.Background(), domain.CategoryEntity{ID: categoryID})

	require.NotNil(t, res)
	require.IsType(t, &domain.CategoryEntity{}, res)
	assert.Equal(t, categoryID, res.ID)
	assert.Equal(t, "name", res.Name.String())
	assert.Equal(t, "description", res.Description.String())

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_ExistsCategory_WhenTheQueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE `categories`.`name` = ? AND `categories`.`description` = ?")).
		WithArgs("name", "category description").
		WillReturnError(fmt.Errorf("some error"))

	repo := NewMySqlCategoriesRepository(db)
	category := domain.CategoryRecord{Name: "name", Description: "category description"}

	res, err := repo.ExistsCategory(context.Background(), category)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_ExistsCategory_WhenItDoesNotFail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `categories` WHERE `categories`.`name` = ? AND `categories`.`description` = ?")).
		WithArgs("name", "category description").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	repo := NewMySqlCategoriesRepository(db)
	category := domain.CategoryRecord{Name: "name", Description: "category description"}

	res, err := repo.ExistsCategory(context.Background(), category)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_GetCategories_WhenItFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)

	repo := NewMySqlCategoriesRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`userId` = ?")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetCategories(context.Background(), domain.CategoryRecord{UserID: 1})

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_GetCategories_When_It_Does_Not_Fail_Including_Items(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)

	repo := NewMySqlCategoriesRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`userId` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(categoryColumns).
			AddRow(11, "category1", "desc 1").
			AddRow(12, "category2", "desc 2"))

	res, err := repo.GetCategories(context.Background(), domain.CategoryRecord{UserID: 1})

	assert.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, "category1", res[0].Name)
	assert.Equal(t, "desc 1", res[0].Description)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, "category2", res[1].Name)
	assert.Equal(t, "desc 2", res[1].Description)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_CreateCategory_When_The_Create_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories` (`name`,`description`,`userId`) VALUES (?,?,?)")).
		WithArgs("name", "category description", 2).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	nvo, _ := domain.NewCategoryNameValueObject("name")
	dvo, _ := domain.NewCategoryDescriptionValueObject("category description")
	category := domain.CategoryEntity{Name: nvo, Description: dvo, UserID: 2}

	repo := NewMySqlCategoriesRepository(db)

	_, err := repo.CreateCategory(context.Background(), &category)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_CreateCategory_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories` (`name`,`description`,`userId`) VALUES (?,?,?)")).
		WithArgs("name", "category description", 2).
		WillReturnResult(sqlmock.NewResult(12, 0))
	mock.ExpectCommit()

	nvo, _ := domain.NewCategoryNameValueObject("name")
	dvo, _ := domain.NewCategoryDescriptionValueObject("category description")
	category := domain.CategoryEntity{Name: nvo, Description: dvo, UserID: 2}

	repo := NewMySqlCategoriesRepository(db)

	res, err := repo.CreateCategory(context.Background(), &category)

	require.NotNil(t, res)
	assert.IsType(t, &domain.CategoryEntity{}, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_DeleteCategory_When_Deleting_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	categoryID := int32(11)

	repo := NewMySqlCategoriesRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(categoryID).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteCategory(context.Background(), domain.CategoryEntity{ID: categoryID})

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_DeleteCategory_When_It_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	categoryID := int32(11)

	repo := NewMySqlCategoriesRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `categories` WHERE `categories`.`id` = ?")).
		WithArgs(categoryID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteCategory(context.Background(), domain.CategoryEntity{ID: categoryID})

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_UpdateCategory_When_The_Update_Fails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `name`=?,`description`=? WHERE `id` = ?")).
		WithArgs("name", "category description", 11).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	repo := NewMySqlCategoriesRepository(db)
	nvo, _ := domain.NewCategoryNameValueObject("name")
	dvo, _ := domain.NewCategoryDescriptionValueObject("category description")
	category := domain.CategoryEntity{ID: 11, Name: nvo, Description: dvo}

	_, err := repo.UpdateCategory(context.Background(), &category)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlCategoriesRepository_UpdateCategory_When_The_Update_Does_Not_Fail(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `name`=?,`description`=? WHERE `id` = ?")).
		WithArgs("name", "category description", 11).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewMySqlCategoriesRepository(db)
	nvo, _ := domain.NewCategoryNameValueObject("name")
	dvo, _ := domain.NewCategoryDescriptionValueObject("category description")
	category := domain.CategoryEntity{ID: 11, Name: nvo, Description: dvo}

	updatedCategory, err := repo.UpdateCategory(context.Background(), &category)

	assert.Nil(t, err)
	assert.IsType(t, &domain.CategoryEntity{}, updatedCategory)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}
