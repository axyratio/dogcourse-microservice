package routes

import (
	"auth-service/controllers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/auth/register", controllers.Register(db)).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", controllers.Login(db)).Methods("POST")
	r.HandleFunc("/api/v1/auth/me", controllers.GetUserProfile(db)).Methods("GET")

	return r
}