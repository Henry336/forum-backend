package main

import (
	"fmt"
	"github.com/heinlinhtet/forum-backend/database"
	"github.com/heinlinhtet/forum-backend/handlers"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// Wrapper to enable CORS for any handler
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. PERMISSION SLIPS (CORS) -- read more about this
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Pre-flight check
		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func main() {
	// 1. Initialize Database. NTS: This will use the connection string from database/db.go
	database.InitDB()

	// 2. Setup Routes
	http.HandleFunc("/login", enableCORS(handlers.LoginHandler))

	http.HandleFunc("/topics", enableCORS(handlers.TopicsHandler))
	http.HandleFunc("/topics/", enableCORS(handlers.TopicsHandler)) // Matches /topics/1/posts

	http.HandleFunc("/posts", enableCORS(handlers.PostsHandler))
	http.HandleFunc("/posts/", enableCORS(handlers.PostsHandler)) // Matches /posts/1 or /posts/1/comments
	// This is to get the PORT from the environment, and to default to 8080 if not found
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server starting on port: " + port + "...")
	http.ListenAndServe(":"+port, nil)
}
