package middleware

import (
	"net/http"
	"strings"

	"booking-service/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth ตรวจสอบ Authorization Header และ validate JWT
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return config.JWT_SECRET, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// เอา user_id และ role_id จาก token ใส่ context
		if uid, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(uid))
		}
		if rid, ok := claims["role_id"].(float64); ok {
			c.Set("role_id", uint(rid))
		}

		c.Next()
	}
}

// UserAuth ตรวจสอบ role_id
func UserAuth(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDVal, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing role"})
			c.Abort()
			return
		}

		roleID, ok := roleIDVal.(uint)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role type"})
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == roleID {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		c.Abort()
	}
}
