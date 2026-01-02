// 1. Define the "Blueprint"
package main

import (
	"database/sql"
	"fmt"
	"log"
)

type HighScore struct {
	Player string
	Score  int
}

type Topic struct {
	Id          int
	Title       string
	Description string
}

func getScores(db *sql.DB) {
	rows, _ := db.Query("SELECT player_name, score FROM high_scores")
	defer rows.Close()

	// 2. Create an empty list (slice) to hold our structs
	var results []HighScore

	for rows.Next() {
		// 3. Create a temporary empty struct for the current row
		var s HighScore

		// 4. Scan directly into the struct's fields using &s.Field
		rows.Scan(&s.Player, &s.Score)

		// 5. Add this completed struct to our list
		results = append(results, s)
	}

	// Now 'results' holds all our data cleanly!
	fmt.Println(results)
}

func getTopics(db *sql.DB) {
	rows, err := db.Query("SELECT id, title, description FROM topics")

	if err != nil {
		log.Fatal("Error encountered while querying:", err)
	}

	defer rows.Close()

	var results []Topic

	for rows.Next() {
		var t Topic

		rows.Scan(&t.Id, &t.Title, &t.Description)

		results = append(results, t)
	}

	fmt.Printf("%+v\n", results)
}
