package main

import (
	"log"
	"net/http"
	"os"

	"auth-service/config"
	"auth-service/internal/models"
	"auth-service/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 🔹 Auto migrate models
	db.AutoMigrate(&models.Role{}, &models.User{})

	// 🔹 Seed role ถ้ายังไม่มี
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count == 0 {
		roles := []models.Role{
			{ID: 1, RoleName: "admin"},
			{ID: 2, RoleName: "user"},
		}
		db.Create(&roles)
	}
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r, db)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
