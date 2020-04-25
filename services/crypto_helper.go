package services

import (
	"golang.org/x/crypto/bcrypt"
)

// CryptoHelper is the interface which contains the methods used to use encrypt the passwords
type CryptoHelper interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

// BcryptHelper is the implementation for CryptoHelper and uses the real bcrypt package
type BcryptHelper struct{}

// NewBcryptHelper returns a new MyBryptProvider
func NewBcryptHelper() *BcryptHelper {
	return new(BcryptHelper)
}

// GenerateFromPassword generates a hashed password
func (b *BcryptHelper) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHashAndPassword checks if the given hashed password and the password matches
func (b *BcryptHelper) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
