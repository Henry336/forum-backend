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

type Post struct {
	Id          int
	TopicId     int
	Title       string
	Description string
	Username    string
	CreatedAt   string
}

type Comment struct {
	Id        int
	PostId    int
	Content   string
	Username  string
	CreatedAt string
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

func createPost(db *sql.DB, p Post) (int, error) {

	sqlStatement := `
	INSERT INTO posts (title, description, topic_id, username) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	// Currently hardcoded
	defaultTopicId := 1
	user := "Henry"

	var id int

	err := db.QueryRow(sqlStatement, p.Title, p.Description, defaultTopicId, user).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func getPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, topic_id, title, description, username, created_at FROM posts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Post

	for rows.Next() {
		var p Post

		err = rows.Scan(&p.Id, &p.TopicId, &p.Title, &p.Description, &p.Username, &p.CreatedAt)

		if err != nil {
			return nil, err
		}

		results = append(results, p)
	}

	return results, nil

}

func getPostById(db *sql.DB, id int) (Post, error) {
	var p Post

	err := db.QueryRow("SELECT id, topic_id, title, description, username, created_at FROM posts WHERE id = $1", id).
		Scan(&p.Id, &p.TopicId, &p.Title, &p.Description, &p.Username, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func deletePost(db *sql.DB, id int) error {
	// Note to self: $1 is used as placeholder to prevent SQL injection (security concerns)
	_, err := db.Exec("DELETE FROM posts WHERE id = $1", id)

	return err
}

func updatePostDescription(db *sql.DB, desc string, id int) error {
	_, err := db.Exec("UPDATE posts SET description = $1 WHERE id = $2", desc, id)

	return err
}

func updatePostTitle(db *sql.DB, titleNew string, id int) error {
	_, err := db.Exec("UPDATE posts SET title = $1 WHERE id = $2", titleNew, id)

	return err
}

func createComment(db *sql.DB, c Comment) error {
	sqlStatement := `
	INSERT INTO comments (content, post_id, username)
	VALUES ($1, $2, $3)	
	`

	_, err := db.Exec(sqlStatement, c.Content, c.PostId, "Henry")

	return err
}

func getCommentsByPostId(db *sql.DB, postId int) ([]Comment, error) {
	rows, err := db.Query("SELECT id, post_id, content FROM comments WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment

		err := rows.Scan(&c.Id, &c.PostId, &c.Content, &c.Username, &c.CreatedAt)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. PERMISSION SLIPS (CORS)
	// Allow any origin or set to "http://localhost:5173" for better security
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Allow the specific methods being used
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")

	// Allow JSON headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Pre-flight check
	if r.Method == "OPTIONS" {
		return
	}

	if strings.Contains(r.URL.Path, "comments") {

		parts := strings.Split(r.URL.Path, "/")

		postIdStr := parts[2]
		postId, err := strconv.Atoi(postIdStr)

		if err != nil {
			http.Error(w, "Invalid Topic ID", 400)
			return
		}

		if r.Method == "GET" {
			comments, err := getCommentsByPostId(db, postId)
			if err != nil {
				http.Error(w, "Database error", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(comments)
			return
		}

		if r.Method == "POST" {
			var c Comment

			err := json.NewDecoder(r.Body).Decode(&c)

			if err != nil {
				http.Error(w, "Invalid JSON", 400)
				return
			}

			c.PostId = postId
			err = createComment(db, c)
			if err != nil {
				fmt.Println("SQL Error:", err)
				http.Error(w, "Database error", 500)
				return
			}
			fmt.Fprintf(w, "Comment added to Post %v", postId)
			return
		}
	}

	if r.Method == "GET" {
		// NTS: Check if the URL has an ID
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Case 1: Just need to fetch a single post
		if len(pathParts) == 2 && pathParts[0] == "posts" {
			id, err := strconv.Atoi(pathParts[1])
			if err != nil {
				http.Error(w, "Invalid ID", 400)
				return
			}

			post, err := getPostById(db, id)
			if err != nil {
				http.Error(w, "Post not found", 404)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
			return
		}

		// Case 2: Fetch ALL posts
		posts, err := getPosts(db)

		if err != nil {
			fmt.Println("SQL Error:", err)
			http.Error(w, "Database error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		return

	} else if r.Method == "POST" {

		var p Post
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		id, err := createPost(db, p)

		if err != nil {
			fmt.Println("Error detected", err)
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully created post with ID: %v", id)
		return

	} else if r.Method == "PATCH" {

		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		var p Post
		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		if p.Title != "" {
			err := updatePostTitle(db, p.Title, id)
			if err != nil {
				http.Error(w, "Database error updating title", 500)
				return
			}
		}

		if p.Description != "" {
			err := updatePostDescription(db, p.Description, id)
			if err != nil {
				http.Error(w, "Database error updating description", 500)
				return
			}
		}

		fmt.Fprintf(w, "Successfully updated Post %v", id)

	} else if r.Method == "DELETE" {
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid ID"+idStr, 400)
			return
		}

		err = deletePost(db, id)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		fmt.Fprintf(w, "Successfully deleted Post %v", id)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
