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

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	// URL is like: /posts/1/comments
	parts := strings.Split(r.URL.Path, "/")
	// NTS: Here, parts[0] = "", parts[1] = "posts", parts[2] = "ID", parts[3] = "comments"

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", 400)
		return
	}

	postIdStr := parts[2]
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", 400)
		return
	}

	// GET comments
	if r.Method == "GET" {
		rows, err := database.DB.Query("SELECT id, post_id, content, username, created_at FROM comments WHERE post_id = $1", postId)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		defer rows.Close()

		var comments []models.Comment
		for rows.Next() {
			var c models.Comment
			err = rows.Scan(&c.Id, &c.PostId, &c.Content, &c.Username, &c.CreatedAt)
			if err != nil {
				continue
			}
			comments = append(comments, c)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
		return
	}

	// CREATE comment
	if r.Method == "POST" {
		var c models.Comment
		err = json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		c.PostId = postId
		sqlStatement := `
		INSERT INTO comments (content, post_id, username)
		VALUES ($1, $2, $3) 
		`
		_, err := database.DB.Exec(sqlStatement, c.Content, c.PostId, c.Username)
		if err != nil {
			fmt.Println("SQL Error:", err)
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Comment added to Post %v", postId)
		return
	}
}
