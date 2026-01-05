package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver import
)

// NTS: Global DB variable so handlers can access it
var DB *sql.DB

func InitDB() {
	var err error
	// This is for checking if there is a cloud database url first
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = "user=postgres dbname=cvwo_forum sslmode=disable password=Heinlinhtet@336"
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database successfully!")
}
