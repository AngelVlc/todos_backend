package models

import "github.com/AngelVlc/todos/dtos"

// List is the model for the list
type List struct {
	ID        int32       `gorm:"type:int(32);primary_key"`
	Name      string      `gorm:"type:varchar(50)"`
	ListItems []*ListItem `gorm:"foreignkey:ListID"`
	UserID    int32       `gorm:"column:userId;type:int(32)"`
}

func (l *List) FromDto(dto *dtos.ListDto) {
	l.Name = dto.Name
}

func (l *List) ToResponseDto() *dtos.ListResponseDto {
	res := dtos.ListResponseDto{
		ID:        l.ID,
		Name:      l.Name,
		ListItems: make([]*dtos.ListItemResponseDto, len(l.ListItems)),
	}

	for i, v := range l.ListItems {
		res.ListItems[i] = &dtos.ListItemResponseDto{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
		}
	}

	return &res
}
