package auth

import (
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("Password can't be empty")
	}

	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Cant craete hash from password: %v", err)
	}

	return hashedPass, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Cant compare password with hash: %v", err)
	}

	return match, nil
}
