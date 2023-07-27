package domain

type CreateUserInput struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	IsAdmin         bool   `json:"isAdmin"`
}
