package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCategoryName_Validates_MinLength(t *testing.T) {
	categoryName, err := NewCategoryNameValueObject("")

	assert.Empty(t, categoryName)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The category name can not be empty", badReqErr.Error())
}

func TestNewCategoryName_Validates_MaxLength(t *testing.T) {
	categoryName, err := NewCategoryNameValueObject("012345678901234567890123456789012345678901234567890")

	assert.Empty(t, categoryName)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The category name can not have more than 12 characters", badReqErr.Error())
}

func TestNewCategoryName_Returns_A_Valid_ItemTitle(t *testing.T) {
	categoryName, err := NewCategoryNameValueObject("a valid name")

	assert.Equal(t, "a valid name", categoryName.String())
	assert.NoError(t, err)
}
