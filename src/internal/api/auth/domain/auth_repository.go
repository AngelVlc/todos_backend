package domain

type AuthRepository interface {
	FindUserByName(userName *AuthUserName) (*AuthUser, error)
	FindUserByID(userID *int32) (*AuthUser, error)
}
