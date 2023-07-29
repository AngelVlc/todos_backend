package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	ID           int32               `gorm:"type:int(32);primary_key"`
	Name         UserNameValueObject `gorm:"type:varchar(10);index:idx_users_name,unique"`
	PasswordHash string              `gorm:"column:passwordHash;type:varchar(100)"`
	IsAdmin      bool                `gorm:"column:isAdmin;type:tinyint"`
}

func (UserRecord) TableName() string {
	return "users"
}

func (u *UserRecord) HasPassword(value UserPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(value))
}

func (u *UserRecord) IsTheAdminUser() bool {
	userNameLowerCase := strings.ToLower(string(u.Name))

	return userNameLowerCase == "admin"
}
