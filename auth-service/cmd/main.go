package main

import (
	"log"
	"net/http"
	"os"

	"auth-service/config"
	"auth-service/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.ConnectDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

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
