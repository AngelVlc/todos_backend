package domain

type ListRecord struct {
	ID         int32             `gorm:"type:int(32);primary_key" json:"id"`
	Name       string            `gorm:"type:varchar(50)" json:"name"`
	UserID     int32             `gorm:"column:userId;type:int(32)" json:"-"`
	ItemsCount int32             `gorm:"column:itemsCount;type:int(32)" json:"itemsCount"`
	Items      []*ListItemRecord `gorm:"foreignKey:ListID" json:"items,omitempty"`
}

func (ListRecord) TableName() string {
	return "lists"
}
