package dtos

// GetUsersResultDto is the struct used as result for the GetUsers method
type GetUsersResultDto struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}
