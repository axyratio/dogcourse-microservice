package routes

import (
	"github.com/gin-gonic/gin"
	"auth-service/internal/controllers"
)

func SetupRouter(handler *controllers.Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handler.RegisterHandler)
		api.POST("/login", handler.LoginHandler)
		api.POST("/logout", handler.LogoutHandler)
		api.GET("/verify", handler.VerifyTokenHandler)
	}

	return r
}