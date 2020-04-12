package main

import (
	"golang.org/x/crypto/bcrypt"
)

func getPasswordHash(p string) (string, error) {
	hasshedPass, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}

	return string(hasshedPass), nil
}
