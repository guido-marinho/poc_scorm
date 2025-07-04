package main

import (
	"github.com/guilherme-gatti/poc_scorm/internal/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":3000")
}