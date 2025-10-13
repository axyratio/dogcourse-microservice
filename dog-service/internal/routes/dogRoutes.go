package routes

import (
	"dog-service/internal/controllers"
	"github.com/gin-gonic/gin"
	"dog-service/internal/middleware"
)

func DogRoutes(r *gin.Engine) {
	r.GET("/dogs", middleware.JWTAuth(), middleware.UserAuth(), controllers.GetAllDogByUserID)
	// r.GET("/dogs/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.GetDogByID)
	r.POST("/dogs", middleware.JWTAuth(),middleware.UserAuth(), controllers.CreateDog)
	r.PATCH("/dogs/:id", middleware.JWTAuth(),middleware.UserAuth(), controllers.UpdateDogByID)
	r.DELETE("/dogs/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.DeleteDogByID)
}
