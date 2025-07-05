package main

import (
	"github.com/guilherme-gatti/poc_scorm/internal/router"
	"github.com/guilherme-gatti/poc_scorm/internal/storage"
)

func main() {
	storage.InitDB("storage/database.db")

	r := router.SetupRouter()
	r.Run(":3000")
}
