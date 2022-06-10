package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateHash(password string) (string, error) {
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), hashErr
}
