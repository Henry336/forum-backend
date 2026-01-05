package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/heinlinhtet/forum-backend/database"
	"github.com/heinlinhtet/forum-backend/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO NOTHING", u.Username)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %v logged in", u.Username)
}
