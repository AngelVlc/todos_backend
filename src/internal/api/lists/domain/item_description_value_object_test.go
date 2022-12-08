package domain

import (
	"bytes"
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemDescription_Validates_MaxLength(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 510; i++ {
		b.WriteString("#")
	}

	itemDescription, err := NewItemDescriptionValueObject(b.String())

	assert.Empty(t, itemDescription)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The item description can not have more than 500 characters", badReqErr.Error())
}

func TestNewItemDescription_Returns_A_Valid_ItemTitle(t *testing.T) {
	itemDescription, err := NewItemDescriptionValueObject("a valid description")

	assert.Equal(t, "a valid description", string(itemDescription))
	assert.NoError(t, err)
}
