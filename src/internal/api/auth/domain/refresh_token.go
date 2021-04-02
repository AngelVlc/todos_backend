package domain

import "time"

type RefreshToken struct {
	ID             int32     `gorm:"type:int(32);primary_key"`
	UserID         int32     `gorm:"column:userId;type:int(32)"`
	RefreshToken   string    `gorm:"column:refreshToken;type:varchar(250);index:idx_refresh_token"`
	ExpirationDate time.Time `gorm:"column:expirationDate;type:timestamp"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
