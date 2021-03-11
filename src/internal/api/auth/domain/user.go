package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int32    `gorm:"type:int(32);primary_key"`
	Name         UserName `gorm:"type:varchar(10);index:idx_users_name"`
	PasswordHash string   `gorm:"type:varchar(100)"`
	IsAdmin      bool     `gorm:"type:tinyint"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) HasPassword(value UserPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(value))
}

func (u *User) IsTheAdminUser() bool {
	userNameLowerCase := strings.ToLower(string(u.Name))
	return userNameLowerCase == "admin"
}
