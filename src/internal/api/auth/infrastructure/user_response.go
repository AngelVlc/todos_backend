package infrastructure

// UserResponse is the struct used to send user info
type UserResponse struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}
