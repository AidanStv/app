package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("secret-key")

func GenerateAccessToken(userID int, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(key)
}

func GenerateRefreshToken(userID int, email string) (string, error) {
	сlaims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		сlaims,
	)
	return token.SignedString(key)
}

func GenerateToken(userID int, email string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"email":   email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return key, nil
	})
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
