// Package db ...
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // sql driver
)

// Connect ... returns a connection to the database
func Connect() *sql.DB {
	username := os.Getenv("PGUSER")
	dbname := os.Getenv("DBNAME")
	dbhost := os.Getenv("PGHOST")
	dbport := os.Getenv("PGPORT")

	pgdb := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", username, dbhost, dbport, dbname)

	store, err := sql.Open("postgres", pgdb)

	if err != nil {
		log.Fatal(err)
	}

	createTodoTable(store)

	return store
}
