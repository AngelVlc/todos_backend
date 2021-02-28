package domain

// JwtClaimsInfo is the struct which contains the jwt token claims
type JwtClaimsInfo struct {
	UserID   int32
	UserName string
	IsAdmin  bool
}
