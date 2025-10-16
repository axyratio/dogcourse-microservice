package routes

import (
	"course-service/internal/controllers"
	"course-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(r *gin.Engine) {
	course := r.Group("/courses")
	course.Use(middleware.AuthMiddleware())
	{
		course.GET("", controllers.GetAllCourses)
		course.GET("/:id", controllers.GetCourseByID)
		course.POST("", middleware.AdminAuth(), controllers.CreateCourse)
		course.PATCH("/:id", middleware.AdminAuth(), controllers.UpdateCourse)
		course.DELETE("/:id", middleware.AdminAuth(), controllers.DeleteCourse)
	}
}


