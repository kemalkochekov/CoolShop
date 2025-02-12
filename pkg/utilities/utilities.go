package utilities

import "golang.org/x/crypto/bcrypt"

const generatePasswordCost = 14

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), generatePasswordCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false
	}

	return true
}
