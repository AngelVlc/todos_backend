package domain

// TokenResponse is the dto used as result for the token endpoint
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
