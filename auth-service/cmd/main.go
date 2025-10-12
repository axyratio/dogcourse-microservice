package main

import (
	"fmt"
	"auth-service/config"
	"auth-service/internal/routes"
	"auth-service/internal/models"
	"auth-service/internal/controllers"
	"auth-service/internal/repositories"

	"os"
)

func main() {
	db := config.ConnectDB()
	db.AutoMigrate(&models.User{}, &models.Role{})
	fmt.Println("✅ Migrated user & role models")

	repo := repositories.NewRepository(db)
	handler := controllers.NewHandler(repo)

	r := routes.SetupRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}