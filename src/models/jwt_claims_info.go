package models

// JwtClaimsInfo is the struct which contains the jwt token claims
type JwtClaimsInfo struct {
	UserID   int32
	UserName string
	IsAdmin  bool
}

// RefreshTokenClaimsInfo is the struct which contains the refresh token claims
type RefreshTokenClaimsInfo struct {
	UserID int32
}
