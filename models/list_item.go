package models

import "github.com/AngelVlc/todos/dtos"

// ListItem is the model for a single list item
type ListItem struct {
	ID          int32  `gorm:"type:int(32);primary_key"`
	ListID      int32  `gorm:"column:listId;type:int(32)"`
	Title       string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(200)"`
}

func (ListItem) TableName() string {
	return "listItems"
}

func (l *ListItem) FromDto(dto *dtos.ListItemDto) {
	l.Title = dto.Title
	l.Description = dto.Description
}
