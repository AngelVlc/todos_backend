package dtos

// UserDto is the struct used as DTO for a user
type UserResponseDto struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}
