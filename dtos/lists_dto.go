package dtos

import "github.com/AngelVlc/todos/models"

// ListDto is the struct used as DTO for a List
type ListDto struct {
	Name string
}

// ToList returns a List from the Dto
func (dto *ListDto) ToList() models.List {
	return models.List{
		Name: dto.Name,
	}
}
