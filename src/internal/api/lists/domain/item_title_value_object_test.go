package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemTitle_Validates_MinLength(t *testing.T) {
	itemTitle, err := NewItemTitleValueObject("")

	assert.Empty(t, itemTitle)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The item title can not be empty", badReqErr.Error())
}

func TestNewItemTitle_Validates_MaxLength(t *testing.T) {
	itemTitle, err := NewItemTitleValueObject("012345678901234567890123456789012345678901234567890")

	assert.Empty(t, itemTitle)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The item title can not have more than 50 characters", badReqErr.Error())
}

func TestNewItemTitle_Returns_A_Valid_ItemTitle(t *testing.T) {
	itemTitle, err := NewItemTitleValueObject("a valid title")

	assert.Equal(t, "a valid title", string(itemTitle))
	assert.NoError(t, err)
}
