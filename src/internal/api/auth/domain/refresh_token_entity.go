package domain

import "time"

type RefreshTokenEntity struct {
	ID             int32
	UserID         int32
	RefreshToken   string
	ExpirationDate time.Time
}

func (e *RefreshTokenEntity) ToRefreshTokenRecord() *RefreshTokenRecord {
	return &RefreshTokenRecord{
		ID:             e.ID,
		UserID:         e.UserID,
		RefreshToken:   e.RefreshToken,
		ExpirationDate: e.ExpirationDate,
	}
}
