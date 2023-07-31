package domain

import (
	"testing"

	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewUserNameValueObject_Validates_MinLength(t *testing.T) {
	userName, err := NewUserNameValueObject("")

	assert.Empty(t, userName)

	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The user name can not be empty", badReqErr.Error())
}

func Test_NewUserNameValueObject_Validates_MaxLength(t *testing.T) {
	userName, err := NewUserNameValueObject("012345678900")

	assert.Empty(t, userName)
	badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
	require.Equal(t, true, isBadReqErr, "should be a bad request error")
	assert.Equal(t, "The user name can not have more than 10 characters", badReqErr.Error())
}

func Test_NewUserNameValueObject_Returns_A_Valid_UserName(t *testing.T) {
	userName, err := NewUserNameValueObject("one name")

	assert.Equal(t, "one name", userName.String())
	assert.NoError(t, err)
}
