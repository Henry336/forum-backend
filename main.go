package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	//createTopic("General Chat", "A place to talk about anything.")
	//createTopic("Help", "Ask for help with the CVWO assignment here, haha.")
	//getTopics()

	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}

func createTopic(db *sql.DB, t Topic) (int, error) {

	sqlStatement := `
	INSERT INTO topics (title, description) 
	VALUES ($1, $2)
	RETURNING id
	`

	var id int

	err := db.QueryRow(sqlStatement, t.Title, t.Description).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
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

func deleteTopic(db *sql.DB, id int) error {
	// Note to self: $1 is used as placeholder to prevent SQL injection (security concerns)
	_, err := db.Exec("DELETE FROM topics WHERE id = $1", id)

	return err
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		topics, err := getTopics(db)

		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topics)

	} else if r.Method == "POST" {

		var t Topic
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		id, err := createTopic(db, t)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully created topic with ID: %v", id)
	} else if r.Method == "DELETE" {
		idStr := strings.TrimPrefix(r.URL.Path, "/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		err = deleteTopic(db, id)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully deleted Topic %v", id)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
