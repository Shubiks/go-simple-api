package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Shubiks/go-simple-api/internal/auth"
	"github.com/Shubiks/go-simple-api/internal/s3"
)

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(auth.ContextUserIDKey)
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "Image not provided", http.StatusBadRequest)
		return
	}

	imageURL, err := s3.UploadProfilePicture(file, fileHeader, userID)
	if err != nil {
		http.Error(w, "Failed to upload to S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE users SET profile_picture_url=$1 WHERE id=$2", imageURL, userID)
	if err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"profile_picture_url": imageURL})
}
