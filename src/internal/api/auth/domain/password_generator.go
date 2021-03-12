package domain

type PasswordGenerator interface {
	GenerateFromPassword(password UserPassword) (string, error)
}
