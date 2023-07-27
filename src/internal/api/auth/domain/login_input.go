package domain

type LoginInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
