package models

import "github.com/AngelVlc/todos/dtos"

type User struct {
	ID           int32  `gorm:"type:int(32);primary_key"`
	Name         string `gorm:"type:varchar(10);index:idx_users_name"`
	PasswordHash string `gorm:"type:varchar(100)"`
	IsAdmin      bool   `gorm:"type:tinyint(100)"`
}

func (u *User) FromDto(dto *dtos.UserDto) {
	u.Name = dto.Name
	u.IsAdmin = dto.IsAdmin
}
