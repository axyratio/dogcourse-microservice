package routes

import (
	"dog-service/internal/controllers"
	"dog-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DogRoutes(r *gin.Engine) {
	r.GET("/dogs", middleware.AuthMiddleware(), controllers.GetAllDogByUserID)
	r.POST("/dogs/batch", middleware.AuthMiddleware(), controllers.GetDogsBatch)
	// r.GET("/dogs/:id", middleware.AuthMiddleware(), controllers.GetDogByID)
	r.POST("/dogs", middleware.AuthMiddleware(), controllers.CreateDog)
	r.PATCH("/dogs/:id", middleware.AuthMiddleware(), controllers.UpdateDogByID)
	r.DELETE("/dogs/:id", middleware.AuthMiddleware(), controllers.DeleteDogByID)
}
