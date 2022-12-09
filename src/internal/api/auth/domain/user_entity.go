package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserEntity struct {
	ID           int32               `gorm:"type:int(32);primary_key"`
	Name         UserNameValueObject `gorm:"type:varchar(10);index:idx_users_name,unique"`
	PasswordHash string              `gorm:"column:passwordHash;type:varchar(100)"`
	IsAdmin      bool                `gorm:"column:isAdmin;type:tinyint"`
}

func (UserEntity) TableName() string {
	return "users"
}

func (u *UserEntity) HasPassword(value UserPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(value))
}

func (u *UserEntity) IsTheAdminUser() bool {
	userNameLowerCase := strings.ToLower(string(u.Name))

	return userNameLowerCase == "admin"
}
