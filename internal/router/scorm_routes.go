package router

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherme-gatti/poc_scorm/internal/scorm"
)

func SetupScormRoutes(r *gin.Engine) {
	r.GET("/progress/:userId", scorm.ProgressHandler)
	r.POST("/track", scorm.TrackHandler)
	r.POST("/upload", scorm.UploadHandler)
	r.GET("/progress/:userId/csv", scorm.ExportCSVHandler)
	r.GET("/progress/:userId/pdf", scorm.ExportPDFHandler)
}
