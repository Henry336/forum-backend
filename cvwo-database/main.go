package main

import (
	"database/sql" // Standard library for SQL interactions
	"fmt"
	"log"

	_ "modernc.org/sqlite" // THE SQL Driver. The '_' is critical (see below)
)

func main() {
	// 1. OPEN CONNECTION
	// This creates 'forum.db' if it doesn't exist
	db, err := sql.Open("sqlite", "forum.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. TEST CONNECTION
	err = db.Ping()
	if err != nil {
		log.Fatal("Database dead:", err)
	}
	fmt.Println("Connected to forum.db!")

	// 3. CREATE TABLE (The Blueprint)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT
	);`

	_, err = db.Exec(createTableSQL) // db.Exec is for commands that don't return rows
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table 'users' ready.")

	// 4. INSERT DATA
	insertSQL := `INSERT INTO users (username, email) VALUES (?, ?)`
	// The '?' are placeholders. Golang injects the values safely.
	_, err = db.Exec(insertSQL, "Hein Lin Htet", "e1682256@u.nus.edu")
	if err != nil {
		log.Fatal("Failed to insert:", err)
	}
	fmt.Println("Hein Lin Htet added to database.")
}
