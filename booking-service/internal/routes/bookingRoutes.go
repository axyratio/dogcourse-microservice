package routes

import (
	"booking-service/internal/controllers"
	"booking-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.Engine) {
	r.GET("courses/booking", middleware.JWTAuth(), middleware.UserAuth(), controllers.GetBookings)
	r.GET("courses/booking/:id", middleware.JWTAuth(), middleware.UserAuth(), controllers.GetBookingByID)
	r.POST("courses/booking/", middleware.JWTAuth(), middleware.UserAuth(), controllers.CreateBooking)
}

func ApproveRoutes(router *gin.Engine) {
    booking := router.Group("/booking")
    {
        booking.PATCH("/:id/approve", middleware.JWTAuth(), middleware.UserAuth(), controllers.ApproveBooking)
        booking.PATCH("/:id/reject", middleware.JWTAuth(), middleware.UserAuth(), controllers.RejectBooking)
    }
}
