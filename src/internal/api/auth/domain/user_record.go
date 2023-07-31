package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	ID           int32  `gorm:"type:int(32);primary_key" json:"id"`
	Name         string `gorm:"type:varchar(10);index:idx_users_name,unique" json:"name"`
	PasswordHash string `gorm:"column:passwordHash;type:varchar(100)" json:"-"`
	IsAdmin      bool   `gorm:"column:isAdmin;type:tinyint" json:"isAdmin"`
}

func (UserRecord) TableName() string {
	return "users"
}

func (u *UserRecord) HasPassword(value string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(value))
}

func (u *UserRecord) IsTheAdminUser() bool {
	userNameLowerCase := strings.ToLower(string(u.Name))

	return userNameLowerCase == "admin"
}
