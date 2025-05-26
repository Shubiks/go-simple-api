package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shubiks/go-simple-api/internal/auth"
	"github.com/go-chi/chi/v5"
)

var followDB *sql.DB

func SetFollowDB(db *sql.DB) {
	followDB = db
}

func SendFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID (the follower) from JWT context
	userIDVal := r.Context().Value(auth.ContextUserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	followerID := userIDVal.(int)

	// Get following user ID from URL
	followingIDStr := chi.URLParam(r, "user_id")
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		http.Error(w, "Invalid following user ID", http.StatusBadRequest)
		return
	}

	if followerID == followingID {
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	fmt.Printf("Sending follow request: followerID=%d, followingID=%d\n", followerID, followingID)
	_, err = followDB.Exec(`INSERT INTO follows (follower_id, following_id, accepted, created_at)
		VALUES ($1, $2, $3, $4)`, followerID, followingID, false, time.Now())
	if err != nil {
		fmt.Print("Failed to insert follow request: %v\n", err)
		http.Error(w, "Failed to send follow request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request sent"})
}
