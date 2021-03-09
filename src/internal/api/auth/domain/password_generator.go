package domain

type PasswordGenerator interface {
	GenerateFromPassword(password *AuthUserPassword) (string, error)
}
