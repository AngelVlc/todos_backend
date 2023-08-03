package domain

type UserRecord struct {
	ID           int32  `gorm:"type:int(32);primary_key" json:"id"`
	Name         string `gorm:"type:varchar(10);index:idx_users_name,unique" json:"name"`
	PasswordHash string `gorm:"column:passwordHash;type:varchar(100)" json:"-"`
	IsAdmin      bool   `gorm:"column:isAdmin;type:tinyint" json:"isAdmin"`
}

func (UserRecord) TableName() string {
	return "users"
}

func (r *UserRecord) ToUserEntity() *UserEntity {
	nvo, _ := NewUserNameValueObject(r.Name)

	return &UserEntity{
		ID:           r.ID,
		Name:         nvo,
		PasswordHash: r.PasswordHash,
		IsAdmin:      r.IsAdmin,
	}
}
