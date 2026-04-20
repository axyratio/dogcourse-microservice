package routes

import (
	"main/controllers"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

// ReviewRoutes sets up the routes for review-related endpoints.
func ReviewRoutes(router *gin.Engine) {
	// Public routes
	router.GET("courses/reviews", controllers.GetReviews)
	router.GET("courses/reviews/:id", controllers.GetReview)

	// User authenticated routes
	router.POST("courses/reviews/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.CreateReview)
	router.PATCH("courses/reviews/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.UpdateReview)

	// Admin authenticated routes
	router.DELETE("courses/reviews/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.DeleteReview)
}
