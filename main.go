package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"net/http"
	"os"
	"workshop-demo/service"
)

func main() {
	db := SetupDB()
	server := service.NewServer(db)
	http.HandleFunc("/", server.ServeHTTP)
	var err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//SetupDB connect database and create a Database object
func SetupDB() *service.Database {
	databaseURL := os.Getenv("CONTACTS_DB_URL")
	if databaseURL == "" {
		panic("CONTACTS_DB_URL must be set!")
	}

	sqlFiles := "./db/migrations"
	if sqlFilesEnv := os.Getenv("CONTACTS_DB_MIGRATIONS"); sqlFilesEnv != "" {
		sqlFiles = sqlFilesEnv
	}
	allErrors, ok := migrate.ResetSync(databaseURL, sqlFiles)
	if !ok {
		panic(fmt.Sprintf("%+v", allErrors))
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to open DB connection: %+v", err))
	}

	return &service.Database{db}
}
