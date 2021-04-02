package passgen

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordGenerator struct{}

func NewBcryptPasswordGenerator() *BcryptPasswordGenerator {
	return &BcryptPasswordGenerator{}
}

func (g *BcryptPasswordGenerator) GenerateFromPassword(password string) (string, error) {
	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hasshedPass), nil
}
