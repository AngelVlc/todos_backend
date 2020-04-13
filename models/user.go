package models

type User struct {
	ID           int32  `gorm:"type:int(32);primary_key"`
	Name         string `gorm:"type:varchar(10);index:idx_users_name"`
	PasswordHash string `gorm:"type:varchar(100)"`
	IsAdmin      bool   `gorm:"type:tinyint(100)"`
}