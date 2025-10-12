package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID uint, role string, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (uint, string, string, error) {
	fmt.Println(tokenString, "tokenString in ValidateJWT")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat := claims["user_id"].(float64)
		role := claims["role"].(string)
		email := claims["email"].(string)
		return uint(userIDFloat), role, email, nil
	}

	return 0, "", "", fmt.Errorf("invalid token")
}