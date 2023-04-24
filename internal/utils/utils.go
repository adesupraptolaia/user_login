package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(salt)
}

func IsPasswordMatch(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
