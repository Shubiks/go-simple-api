package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Shubiks/go-simple-api/internal/utils"
	"github.com/Shubiks/go-simple-api/models"
	"github.com/jmoiron/sqlx"
)

func GetUsersHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		err := db.Select(&users, "SELECT * FROM users")
		if err != nil {
			log.Printf("Error selecting users: %v", err)
			http.Error(w, "Could not fetch users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

var db *sqlx.DB

func SetDB(database *sqlx.DB) {
	db = database
}

// POST /users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", input.Name, input.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
