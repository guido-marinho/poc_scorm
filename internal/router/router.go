package router

import (
	"github.com/guilherme-gatti/poc_scorm/internal/scorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", scorm.PingHandler)

	// rota para servir o index.html do pacote SCORM
	r.Static("/packages", "./storage/")

	SetupScormRoutes(r)

	return r
}
