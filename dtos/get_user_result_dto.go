package dtos

// GetUserResultDto is the struct used as result for the GetUsers method
type GetUserResultDto struct {
	ID      int32  `json:"id" gorm:"type:int(32);primary_key"`
	Name    string `json:"name" gorm:"type:varchar(10)"`
	IsAdmin bool   `json:"isAdmin" gorm:"type:tinyint(100)"`
}

func (GetUserResultDto) TableName() string {
	return "users"
}