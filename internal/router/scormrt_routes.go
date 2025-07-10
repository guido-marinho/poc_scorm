package router

import (
	"github.com/gin-gonic/gin"
	scormrt "github.com/guilherme-gatti/poc_scorm/internal/scormrt"
)

func SetupScormrtRoutes(r *gin.Engine) {
	r.POST("/scormrt", scormrt.RuntimeHandler)
}
