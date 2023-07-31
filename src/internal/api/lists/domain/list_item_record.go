package domain

type ListItemRecord struct {
	ID          int32  `gorm:"type:int(32);primary_key" json:"id"`
	ListID      int32  `gorm:"column:listId;type:int(32)" json:"-"`
	UserID      int32  `gorm:"column:userId;type:int(32)" json:"-"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(200)" json:"description"`
	Position    int32  `gorm:"column:position;type:int(32)" json:"position"`
}

func (ListItemRecord) TableName() string {
	return "listItems"
}
