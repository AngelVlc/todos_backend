package domain

// TokenClaimsInfo is the struct which contains the jwt token claims
type TokenClaimsInfo struct {
	UserID   int32
	UserName string
	IsAdmin  bool
}
