package util

import "golang.org/x/crypto/bcrypt"

func PasswordCompare(hashed string, password []byte) (bool, error) {
	hashByte := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(hashByte, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func PasswordHash(password string) (string, error) {
	pwByte := []byte(password)
	hashed, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
