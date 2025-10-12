package main

import (
	"course-service/config"
	"course-service/internal/models"
	"course-service/internal/routes"
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
	db.AutoMigrate(&models.Course{})

	r := gin.Default()
	routes.CourseRoutes(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}