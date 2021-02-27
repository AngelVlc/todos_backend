package dtos

// UserDto is the struct used as DTO for a user
type UserDto struct {
	Name               string `json:"name"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
	IsAdmin            bool   `json:"isAdmin"`
}
