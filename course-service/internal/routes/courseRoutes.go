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
		course.GET("", middleware.UserAuth(), controllers.GetAllCourses)
		course.GET("/:id", middleware.UserAuth(), controllers.GetCourseByID)
		course.POST("", middleware.UserAuth(), controllers.CreateCourse)
		course.PATCH("/:id", middleware.UserAuth(), controllers.UpdateCourse)
		course.DELETE("/:id", middleware.UserAuth(),  controllers.DeleteCourse)
	}
}
