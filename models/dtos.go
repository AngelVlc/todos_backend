package models

// UserDto is the struct used as DTO for a user
type UserDto struct {
	Name               string
	NewPassword        string
	ConfirmNewPassword string
	IsAdmin            bool
}

// ToUser returns a User from the Dto
func (dto *UserDto) ToUser() User {
	return User{
		Name:    dto.Name,
		IsAdmin: dto.IsAdmin,
	}
}

// GetUsersResultDto is the struct used as result for the GetUsers method
type GetUsersResultDto struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}
