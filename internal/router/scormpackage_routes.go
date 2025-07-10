package router

import (
	"github.com/gin-gonic/gin"
	scorm "github.com/guilherme-gatti/poc_scorm/internal/scormpackage"
)

func SetupScormPackageRoutes(r *gin.Engine) {

	r.GET("/progress/:userId", scorm.ProgressHandler)
	r.POST("/track", scorm.TrackHandler)
	r.POST("/upload", scorm.UploadHandler)
	r.GET("/progress/:userId/csv", scorm.ExportCSVHandler)
	r.GET("/progress/:userId/pdf", scorm.ExportPDFHandler)

	r.GET("/courses", scorm.ListCoursesHandler)
	r.GET("/courses/:id/validated", scorm.GetCourseValidatedHandler)
	r.GET("/courses/:id/view", scorm.GetCourseValidatedHandler)
	r.POST("/courses/:id/validate", scorm.ValidateExistingCourseHandler)
	r.DELETE("/courses/:id", scorm.DeleteCourseHandler)
}
