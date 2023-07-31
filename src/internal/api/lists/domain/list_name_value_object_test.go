package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewListName_Validates_MinLength(t *testing.T) {
	listName, err := NewListNameValueObject("")

	assert.Empty(t, listName)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The list name can not be empty", badReqErr.Error())
}

func TestNewListName_Validates_MaxLength(t *testing.T) {
	listName, err := NewListNameValueObject("012345678901234567890123456789012345678901234567890")

	assert.Empty(t, listName)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The list name can not have more than 50 characters", badReqErr.Error())
}

func TestNewListName_Returns_A_Valid_ListName(t *testing.T) {
	listName, err := NewListNameValueObject("a valid name")

	assert.Equal(t, "a valid name", listName.String())
	assert.NoError(t, err)
}
