package dtos

import "github.com/AngelVlc/todos/models"

// UserDto is the struct used as DTO for a user
type UserDto struct {
	Name               string `json:"name"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
	IsAdmin            bool   `json:"isAdmin"`
}

// ToUser returns a User from the Dto
func (dto *UserDto) ToUser() models.User {
	return models.User{
		Name:    dto.Name,
		IsAdmin: dto.IsAdmin,
	}
}
