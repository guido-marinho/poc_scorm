package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSource string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	// executa o schema na inicializacao
	schema, err := os.ReadFile("internal/storage/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Banco de dados inicializado com sucesso")
}
