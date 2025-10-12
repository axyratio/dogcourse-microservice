package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetUserIDFromRequest(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := ValidateJWT(token)
	if err != nil {
		return 0, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid token")
	}

	return uint(userID), nil
}
