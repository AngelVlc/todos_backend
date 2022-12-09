package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewListName_Validates_MinLength(t *testing.T) {
	listName, err := NewUserNameValueObject("")

	assert.Empty(t, listName)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The user name can not be empty", badReqErr.Error())
}

func TestNewListName_Validates_MaxLength(t *testing.T) {
	listName, err := NewUserNameValueObject("012345678900")

	assert.Empty(t, listName)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The user name can not have more than 10 characters", badReqErr.Error())
}

func TestNewListName_Returns_A_Valid_ItemTitle(t *testing.T) {
	listName, err := NewUserNameValueObject("one name")

	assert.Equal(t, "one name", string(listName))
	assert.NoError(t, err)
}
