package domain

import (
	"bytes"
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCategoryDescription_Validates_MaxLength(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 510; i++ {
		b.WriteString("#")
	}

	categoryDescription, err := NewCategoryDescriptionValueObject(b.String())

	assert.Empty(t, categoryDescription)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The category description can not have more than 500 characters", badReqErr.Error())
}

func TestNewCategoryDescription_Returns_A_Valid_ItemDescription(t *testing.T) {
	categoryDescription, err := NewCategoryDescriptionValueObject("a valid description")

	assert.Equal(t, "a valid description", categoryDescription.String())
	assert.NoError(t, err)
}
