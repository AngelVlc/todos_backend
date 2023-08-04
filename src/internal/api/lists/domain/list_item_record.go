package domain

type ListItemRecord struct {
	ID          int32  `gorm:"type:int(32);primary_key"`
	ListID      int32  `gorm:"column:listId;type:int(32)"`
	UserID      int32  `gorm:"column:userId;type:int(32)"`
	Title       string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(200)"`
	Position    int32  `gorm:"column:position;type:int(32)"`
}

func (ListItemRecord) TableName() string {
	return "listItems"
}
