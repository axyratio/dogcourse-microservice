package utils
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, error) {
userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "เป็นคนไม่มีสิทธิ์"})
		return 0, nil
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return 0, nil
	}

	return userID, nil
}