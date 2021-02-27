package dtos

// TokenDto is the dto used for login
type TokenDto struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
