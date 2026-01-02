package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Topic struct {
	Id          int
	Title       string
	Description string
}

func main() {
	// 1. Connection settings
	connStr := "user=postgres password=Heinlinhtet@336 dbname=cvwo_forum sslmode=disable"

	// 2. Open the connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Test the connection (Ping)
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	fmt.Println("Connected successfully!")

	//createTopic(db, "General Chat", "A place to talk about anything.")
	//createTopic(db, "Help", "Ask for help with the CVWO assignment here, haha.")
	//getTopics(db)

	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}

func createTopic(db *sql.DB, title string, description string) {
	sqlStatement := `
	INSERT INTO topics (title, description) 
	VALUES ($1, $2)
	RETURNING id
	`

	var id int

	err := db.QueryRow(sqlStatement, title, description).Scan(&id)

	if err != nil {
		log.Fatal("Failed to insert text:", err)
	}

	fmt.Printf("Created Topic with Title: %v, and ID: %v\n", title, id)
}

func getTopics(db *sql.DB) ([]Topic, error) {
	rows, err := db.Query("SELECT id, title, description FROM topics")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Topic

	for rows.Next() {
		var t Topic

		rows.Scan(&t.Id, &t.Title, &t.Description)

		results = append(results, t)
	}

	return results, nil

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	topics, err := getTopics(db)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
