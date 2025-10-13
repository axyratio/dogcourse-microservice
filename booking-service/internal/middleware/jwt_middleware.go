package middleware

import (
	"booking-service/internal/utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println(authHeader, "authheader")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		userID, role, email, err := utils.VerifyToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}


		// ✅ เก็บข้อมูลไว้ใน context
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("email", email)

		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtKey = []byte(os.Getenv("JWT_SECRET"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// ✅ ดึง user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload: missing user_id"})
			return
		}
		userID := uint(userIDFloat)

		// ✅ ดึง role (string) จาก token
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload: missing role"})
			return
		}

		// ✨ เก็บลง context
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}


func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}







func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		fmt.Println(role, "role in user auth middleware")
		if !exists || role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}