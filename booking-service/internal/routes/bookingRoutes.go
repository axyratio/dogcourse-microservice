package routes

import (
	"booking-service/internal/controllers"
	"booking-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.Engine) {
	r.GET("courses/booking", middleware.AuthMiddleware(), controllers.GetBookings)
	r.GET("courses/booking/:id", middleware.AuthMiddleware(), controllers.GetBookingByID)
	r.GET("courses/booked/:id", middleware.AuthMiddleware(), controllers.GetBooked) 
	r.POST("courses/booking/", middleware.AuthMiddleware(), controllers.CreateBooking)

}

func ApproveRoutes(router *gin.Engine) {
    booking := router.Group("/booking")
    {
        booking.PATCH("/:id/approve", middleware.AuthMiddleware(), middleware.AdminAuth(), controllers.ApproveBooking)
        booking.PATCH("/:id/reject", middleware.AuthMiddleware(),  middleware.AdminAuth(), controllers.RejectBooking)
    }
}
