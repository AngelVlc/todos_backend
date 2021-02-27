package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHelper(t *testing.T) {
	prv := NewBcryptHelper()

	password := "the_password"

	hashedBytes, err := prv.GenerateFromPassword([]byte(password))

	assert.Nil(t, err)

	err = prv.CompareHashAndPassword(hashedBytes, []byte(password))

	assert.Nil(t, err)
}
