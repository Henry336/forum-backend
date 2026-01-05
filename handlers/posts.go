package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/heinlinhtet/forum-backend/database"
	"github.com/heinlinhtet/forum-backend/models"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	// If the URL contains "comments", hand it off to the Comments Handler
	if strings.Contains(r.URL.Path, "/comments") {
		CommentsHandler(w, r)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Case 1: Get only one post (/posts/{id})
	if r.Method == "GET" && len(pathParts) == 2 {
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		var p models.Post
		err = database.DB.QueryRow("SELECT id, topic_id, title, description, username, created_at FROM posts WHERE id = $1", id).
			Scan(&p.Id, &p.TopicId, &p.Title, &p.Description, &p.Username, &p.CreatedAt)
		if err != nil {
			http.Error(w, "Post not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
		return
	}

	// Case 2: Get all posts (/posts)
	if r.Method == "GET" {
		rows, err := database.DB.Query("SELECT id, topic_id, title, description, username, created_at FROM posts")
		if err != nil {
			fmt.Println("SQL Error:", err)
			http.Error(w, "Database error", 500)
			return
		}
		defer rows.Close()

		var results []models.Post
		for rows.Next() {
			var p models.Post
			err = rows.Scan(&p.Id, &p.TopicId, &p.Title, &p.Description, &p.Username, &p.CreatedAt)
			if err != nil {
				continue
			}
			results = append(results, p)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	// Case 3: Create a post (POST /posts)
	if r.Method == "POST" {
		var p models.Post
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		sqlStatement := `
		INSERT INTO posts (title, description, topic_id, username) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`
		var id int
		err = database.DB.QueryRow(sqlStatement, p.Title, p.Description, p.TopicId, p.Username).Scan(&id)
		if err != nil {
			fmt.Println("Error detected", err)
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully created post with ID: %v", id)
		return
	}

	// Case 4: Update a post (PATCH /posts/{id})
	if r.Method == "PATCH" {
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		var p models.Post
		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		if p.Title != "" {
			_, err := database.DB.Exec("UPDATE posts SET title = $1 WHERE id = $2", p.Title, id)
			if err != nil {
				http.Error(w, "Database error updating title", 500)
				return
			}
		}

		if p.Description != "" {
			_, err := database.DB.Exec("UPDATE posts SET description = $1 WHERE id = $2", p.Description, id)
			if err != nil {
				http.Error(w, "Database error updating description", 500)
				return
			}
		}

		fmt.Fprintf(w, "Successfully updated Post %v", id)
		return
	}

	// Case 5: Delete a post (DELETE /posts/{id})
	if r.Method == "DELETE" {
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		_, err = database.DB.Exec("DELETE FROM posts WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully deleted Post %v", id)
		return
	}
}
