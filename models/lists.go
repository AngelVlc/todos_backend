package models

// List is the model for the list
type List struct {
	ID        int32       `gorm:"type:int(32);primary_key"`
	Name      string      `gorm:"type:varchar(50)"`
	ListItems []*ListItem `gorm:"foreignkey:ListID"`
	UserID    int32       `gorm:"column:userId;type:int(32)"`
}
