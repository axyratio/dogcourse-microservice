package routes

import (
	"main/controllers"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(r *gin.Engine) {
	r.GET("/courses", middleware.JWTAuth(), middleware.UserAuth(),controllers.GetAllCourses)
	r.GET("/courses/:id", middleware.JWTAuth(), middleware.UserAuth(),controllers.GetCourseByID)
	r.POST("/courses", middleware.JWTAuth(), middleware.UserAuth(),controllers.CreateCourse)
	r.PATCH("/courses/:id", middleware.JWTAuth(), middleware.UserAuth(),controllers.UpdateCourse)
	r.DELETE("/courses/:id", middleware.JWTAuth(),middleware.UserAuth(), controllers.DeleteCourse)
}
