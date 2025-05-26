package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Gets the jwt secret from the env file
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// Creates jwt token for user authentication
func CreateToken(userID, email string) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
