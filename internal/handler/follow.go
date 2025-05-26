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
		fmt.Printf("Failed to insert follow request: %v\n", err)
		http.Error(w, "Failed to send follow request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request sent"})
}

func AcceptFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Extract currently logged-in user ID (this user is the one being followed)
	userIDVal := r.Context().Value(auth.ContextUserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	followingID := userIDVal.(int)

	// Extract follower's ID from the URL param
	followerIDStr := chi.URLParam(r, "user_id")
	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Check if a follow request exists
	var exists bool
	err = followDB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2 AND accepted = false
		)
	`, followerID, followingID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "No pending follow request found", http.StatusNotFound)
		return
	}

	// Update the follow request to accepted
	_, err = followDB.Exec(`
		UPDATE follows
		SET accepted = true
		WHERE follower_id = $1 AND following_id = $2
	`, followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Accepted follow request from user %d to user %d\n", followerID, followingID)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request accepted",
	})
}
