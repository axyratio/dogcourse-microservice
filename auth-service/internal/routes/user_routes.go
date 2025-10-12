package routes

import (
	"auth-service/internal/controllers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *mux.Router, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	r.HandleFunc("/api/v1/auth/register", authController.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authController.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/profile", authController.Profile).Methods("GET")
}
