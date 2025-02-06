package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("problem hashing password: %w", err)
	}

	return string(hashedPassword), nil

}

func ComparePassword(hashedPassword string, passwordRequest string) error  {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordRequest))

}