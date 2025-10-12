package main

import (
	"booking-service/config"
	"booking-service/models"
	"booking-service/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	
	config.LoadConfig()
	config.ConnectDB()

	// สร้างตารางอัตโนมัติ
	config.DB.AutoMigrate(&models.Booking{}, &models.BookingDog{})

	r := gin.Default()

	// ติดตั้ง routes
	routes.BookingRoutes(r)
	routes.ApproveRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Booking service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
