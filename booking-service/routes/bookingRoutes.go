package routes

import (
	"booking-service/controllers"
	"booking-service/middleware"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.Engine) {
	booking := r.Group("/api/v1/bookings")
		booking.Use(middleware.JWTAuth())  // middleware JWTAuth
	{
		booking.GET("", controllers.GetBookings)      // <-- ไม่มี "/"
		booking.GET("/:id", controllers.GetBookingByID)
		booking.POST("", controllers.CreateBooking)   // <-- ไม่มี "/"
	}
}
func ApproveRoutes(r *gin.Engine) {
	approve := r.Group("/api/v1/bookings")
	{
		approve.PATCH("/:id/approve", middleware.JWTAuth(), controllers.ApproveBooking)
		approve.PATCH("/:id/reject", middleware.JWTAuth(), controllers.RejectBooking)
	}
}
