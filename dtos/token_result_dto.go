package dtos

// TokenResultDto is the dto used as result for the token endpoint
type TokenResultDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
