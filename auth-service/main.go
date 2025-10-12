package main

import (
	"auth-service/routes"
	"log"
	"net/http"
	"os"
	"auth-service/models"
	"auth-service/config"
	"github.com/joho/godotenv"
)

func main() {


	
	// โหลด env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDB()

	// สร้างตารางอัตโนมัติ
	config.DB.AutoMigrate(&models.Role{}, &models.User{})

	// Seed roles ถ้ายังไม่มี
	var count int64
	config.DB.Model(&models.Role{}).Count(&count)
	if count == 0 {
		roles := []models.Role{
			{ID: 1, RoleName: "admin"},
			{ID: 2, RoleName: "user"},
		}
		if err := config.DB.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed roles: %v", err)
		}
		log.Println("Seeded roles successfully")
	}

	r := routes.SetupRouter(config.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
