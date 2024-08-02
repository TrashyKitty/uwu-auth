package utils

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func IsCorrectPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

func MakeToken(userID string) string {
	secretKey := []byte("very secret key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		userID: userID,
	})

	signedToken, err := token.SignedString(secretKey)

	if err != nil {
		return ""
	}

	return signedToken
}
