package main

import (
	"dog-service/config"
	"dog-service/internal/models"
	"dog-service/internal/routes"
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
	db.AutoMigrate(&models.Dog{})

	r := gin.Default()
	routes.DogRoutes(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(":" + port)
}