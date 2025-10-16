package routes

import (
	"review-service/internal/controllers"
	"review-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ReviewRoutes sets up the routes for review-related endpoints.
func ReviewRoutes(router *gin.Engine) {

	// Public routes
	router.GET("courses/reviews", controllers.GetReviews)
	router.GET("courses/reviews/:id", controllers.GetReview)

	// User authenticated routes
	router.POST("courses/reviews/:id",middleware.AuthMiddleware(), controllers.CreateReview)
	router.PATCH("courses/reviews/:id",middleware.AuthMiddleware(), controllers.UpdateReview)

	// Admin authenticated routes
	router.DELETE("courses/reviews/:id",middleware.AuthMiddleware(), controllers.DeleteReview)
}
