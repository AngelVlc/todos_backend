package infrastructure

type LoginResponse struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	UserID       int32  `json:"userId"`
	UserName     string `json:"userName"`
	IsAdmin      bool   `json:"isAdmin"`
}
