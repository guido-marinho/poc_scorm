package router

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherme-gatti/poc_scorm/internal/course"
)

func SetupCourseRoutes(r *gin.Engine) {
	r.DELETE("/courses/:id", course.DeleteCourseHandler)
	r.GET("/courses", course.CoursesHandler)
}
