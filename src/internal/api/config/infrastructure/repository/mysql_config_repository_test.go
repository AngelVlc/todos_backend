//go:build !e2e
// +build !e2e

package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/AngelVlc/todos_backend/internal/api/config/domain"
	"github.com/AngelVlc/todos_backend/internal/api/shared/infrastructure/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	allowedOriginsColumns = []string{"id", "origin"}
)

func TestMySqlConfigRepository_ExistsAllowedOrigin_QueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	origin := domain.Origin("one_origin")

	expectedExistsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `allowed_origins` WHERE `allowed_origins`.`origin` = ?")).
			WithArgs("one_origin")
	}

	expectedExistsQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.ExistsAllowedOrigin(context.Background(), origin)

	assert.False(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_ExistsAllowedOrigin_Ok(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	origin := domain.Origin("one_origin")

	expectedExistsQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `allowed_origins` WHERE `allowed_origins`.`origin` = ?")).
			WithArgs("one_origin")
	}

	expectedExistsQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	res, err := repo.ExistsAllowedOrigin(context.Background(), origin)

	assert.True(t, res)
	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_GetAllAllowedOrigins_QueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	expectedGetUsersQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,origin FROM `allowed_origins`"))
	}

	expectedGetUsersQuery().WillReturnError(fmt.Errorf("some error"))

	res, err := repo.GetAllAllowedOrigins(context.Background())

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_GetAllAllowedOrigins_Ok(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	expectedGetUsersQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery(regexp.QuoteMeta("SELECT id,origin FROM `allowed_origins`"))
	}

	expectedGetUsersQuery().WillReturnRows(sqlmock.NewRows(allowedOriginsColumns).AddRow(11, "origin1").AddRow(12, "origin2"))

	res, err := repo.GetAllAllowedOrigins(context.Background())

	assert.Nil(t, err)
	require.Equal(t, 2, len(res))
	assert.Equal(t, int32(11), res[0].ID)
	assert.Equal(t, domain.Origin("origin1"), res[0].Origin)
	assert.Equal(t, int32(12), res[1].ID)
	assert.Equal(t, domain.Origin("origin2"), res[1].Origin)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_CreateAllowedOrigin_QueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	allowedOrigin := domain.AllowedOrigin{Origin: domain.Origin("one origin")}

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `allowed_origins` (`origin`) VALUES (?)")).
			WithArgs(allowedOrigin.Origin)
	}

	mock.ExpectBegin()
	expectedInsertExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.CreateAllowedOrigin(context.Background(), &allowedOrigin)

	assert.EqualError(t, err, "some error")

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_CreateAllowedOrigin_Ok(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	allowedOrigin := domain.AllowedOrigin{Origin: domain.Origin("one origin")}

	expectedInsertExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `allowed_origins` (`origin`) VALUES (?)")).
			WithArgs(allowedOrigin.Origin)
	}

	result := sqlmock.NewResult(12, 1)

	mock.ExpectBegin()
	expectedInsertExec().WillReturnResult(result)
	mock.ExpectCommit()

	err := repo.CreateAllowedOrigin(context.Background(), &allowedOrigin)

	assert.Nil(t, err)

	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_DeleteAllowedOrigin_QueryFails(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `allowed_origins` WHERE `allowed_origins`.`id` = ?")).
			WithArgs(1)
	}

	id := int32(1)

	mock.ExpectBegin()
	expectedDeleteExec().WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := repo.DeleteAllowedOrigin(context.Background(), id)

	assert.EqualError(t, err, "some error")
	helpers.CheckSqlMockExpectations(mock, t)
}

func TestMySqlConfigRepository_DeleteAllowedOrigin_Ok(t *testing.T) {
	mock, db := helpers.GetMockedDb(t)
	repo := NewMySqlConfigRepository(db)

	expectedDeleteExec := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `allowed_origins` WHERE `allowed_origins`.`id` = ?")).
			WithArgs(1)
	}

	id := int32(1)

	mock.ExpectBegin()
	expectedDeleteExec().WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteAllowedOrigin(context.Background(), id)

	assert.Nil(t, err)
	helpers.CheckSqlMockExpectations(mock, t)
}
