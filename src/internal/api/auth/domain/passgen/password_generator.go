package passgen

type PasswordGenerator interface {
	GenerateFromPassword(password string) (string, error)
}
