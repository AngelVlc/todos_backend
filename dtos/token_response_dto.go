package dtos

// TokenResponseDto is the dto used as result for the token endpoint
type TokenResponseDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
