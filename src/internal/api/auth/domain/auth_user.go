package domain

import "golang.org/x/crypto/bcrypt"

type AuthUser struct {
	ID           int32        `gorm:"type:int(32);primary_key"`
	Name         AuthUserName `gorm:"type:varchar(10);index:idx_users_name"`
	PasswordHash string       `gorm:"type:varchar(100)"`
	IsAdmin      bool         `gorm:"type:tinyint"`
}

func (AuthUser) TableName() string {
	return "users"
}

func (u *AuthUser) HasPassword(value AuthUserPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(value))
}
