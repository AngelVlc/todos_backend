package models

// Login is the model used for login
type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// RefreshToken is the model used for refreshing the token
type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}
