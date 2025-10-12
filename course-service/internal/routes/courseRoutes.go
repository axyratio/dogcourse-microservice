package routes

import (
	"course-service/internal/controllers"
	"course-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(r *gin.Engine) {
	course := r.Group("/courses")
	course.Use(middleware.JWTAuth())
	{
		course.GET("", controllers.GetAllCourses)
		course.GET("/:id", controllers.GetCourseByID)
		course.POST("", controllers.CreateCourse)
		course.PATCH("/:id", controllers.UpdateCourse)
		course.DELETE("/:id", controllers.DeleteCourse)
	}
}
