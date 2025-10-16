package routes

import (
	"auth-service/internal/controllers"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *controllers.Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handler.RegisterHandler)
		api.POST("/login", handler.LoginHandler)
		api.POST("/logout", handler.LogoutHandler)
		api.GET("/verify", handler.VerifyTokenHandler)
		api.GET("/me", middleware.AuthMiddleware(), handler.GetUserByIDHandler)
	}

	return r
}