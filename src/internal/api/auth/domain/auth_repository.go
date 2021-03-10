package domain

type AuthRepository interface {
	FindUserByName(userName *AuthUserName) (*AuthUser, error)
	FindUserByID(userID *int32) (*AuthUser, error)
	GetAllUsers() ([]*AuthUser, error)
	CreateUser(user *AuthUser) (int32, error)
	DeleteUser(userID *int32) error
	UpdateUser(user *AuthUser) error
}
