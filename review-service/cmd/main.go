package main

import (
	"review-service/config"
	"review-service/internal/models"
	"review-service/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.ConnectDB()
	db.AutoMigrate(&models.Review{})

	r := gin.Default()
	routes.ReviewRoutes(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	r.Run(":" + port)
}