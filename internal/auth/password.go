package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const minLength = 1

func HashPassword(password string) (string, error) {
	if len(password) < minLength {
		return "", fmt.Errorf("password required to be %d characters long", minLength)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
