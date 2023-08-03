package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserEntity struct {
	ID           int32
	Name         UserNameValueObject
	PasswordHash string
	IsAdmin      bool
}

func (e *UserEntity) ToUserRecord() *UserRecord {
	return &UserRecord{
		ID:           e.ID,
		Name:         e.Name.String(),
		PasswordHash: e.PasswordHash,
		IsAdmin:      e.IsAdmin,
	}
}

func (e *UserEntity) HasPassword(value string) error {
	return bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(value))
}

func (e *UserEntity) IsTheAdminUser() bool {
	userNameLowerCase := strings.ToLower(e.Name.String())

	return userNameLowerCase == "admin"
}
