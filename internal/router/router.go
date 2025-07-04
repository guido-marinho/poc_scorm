package router

import (
	"github.com/guilherme-gatti/poc_scorm/internal/scorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// roteador principal - (app do express)
	r := gin.Default()

	// rota get para teste
	// define q c é ponteiro para o context (req + res do express)
	r.GET("/ping", scorm.PingHandler)
	r.POST("/upload", scorm.UploadHandler)

	return r
}
