package domain

type List struct {
	ID          int32    `gorm:"type:int(32);primary_key"`
	Name        ListName `gorm:"type:varchar(50)"`
	UserID      int32    `gorm:"column:userId;type:int(32)"`
	ItemsCount  int32    `gorm:"column:itemsCount;type:int(32)"`
	IsQuickList bool     `gorm:"column:isQuickList;type:tinyint"`
}

func (List) TableName() string {
	return "lists"
}
