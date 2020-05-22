package dtos

import "github.com/AngelVlc/todos/models"

// ListItemDto is the struct used as DTO for a ListItem
type ListItemDto struct {
	Title       string
	Description string
}

// ToList returns a ListItem from the Dto
func (dto *ListItemDto) ToListItem() models.ListItem {
	return models.ListItem{
		Title:       dto.Title,
		Description: dto.Description,
	}
}
