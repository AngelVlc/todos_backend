package domain

type AuthRepository interface {
	FindUserByName(userName *UserName) (*User, error)
	FindUserByID(userID *int32) (*User, error)
	GetAllUsers() ([]*User, error)
	CreateUser(user *User) (int32, error)
	DeleteUser(userID *int32) error
	UpdateUser(user *User) error
}
