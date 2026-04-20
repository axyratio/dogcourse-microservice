package utils

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromRequest(r *http.Request) (uint, error) {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		return 0, errors.New("no token")
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat := claims["user_id"].(float64)
		return uint(idFloat), nil
	}

	return 0, errors.New("invalid token")
}
