package main

import (
	"booking-service/config"
	"booking-service/internal/models"
	"booking-service/internal/routes"
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
	db.AutoMigrate(&models.Booking{})

	r := gin.Default()
	routes.BookingRoutes(r)

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