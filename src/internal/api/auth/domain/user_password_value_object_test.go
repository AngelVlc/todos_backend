package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewUserPasswordValueObject_Validates_MinLength(t *testing.T) {
	userPassword, err := NewUserPasswordValueObject("")

	assert.Empty(t, userPassword)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "Password can not be empty", badReqErr.Error())
}

func Test_NewUserPasswordValueObject_Returns_A_Valid_UserName(t *testing.T) {
	userPassword, err := NewUserPasswordValueObject("one password")

	assert.Equal(t, "one password", userPassword.String())
	assert.NoError(t, err)
}
