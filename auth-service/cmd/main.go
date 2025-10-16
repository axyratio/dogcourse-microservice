package main

import (
	"auth-service/config"
	"auth-service/internal/controllers"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/routes"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
)


func init() {
    if err := godotenv.Load(); err != nil {
        // ถ้า deploy บน Docker/K8s ไม่มีไฟล์ .env ก็โอเค แต่ในเครื่องต้องมี
        log.Println("WARN: .env not loaded:", err)
    }
}

func main() {
	db := config.ConnectDB()
	db.AutoMigrate(&models.User{}, &models.Role{})
	fmt.Println("✅ Migrated user & role models")
	config.SeedRoles(db)

	repo := repositories.NewRepository(db)
	handler := controllers.NewHandler(repo)

	r := routes.SetupRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}