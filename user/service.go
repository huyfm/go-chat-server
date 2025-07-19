package user

import "golang.org/x/crypto/bcrypt"

func HashPassword(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed), err
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
