package router

import (
	scorm "github.com/guilherme-gatti/poc_scorm/internal/scormpackage"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", scorm.PingHandler)

	// rota para servir o index.html do pacote SCORM
	r.Static("/packages", "./storage/")

	SetupScormPackageRoutes(r)
	SetupScormrtRoutes(r)

	return r
}
