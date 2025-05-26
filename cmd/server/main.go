package main

import (
	"log"
	"net/http"

	"github.com/Shubiks/go-simple-api/internal/auth"
	"github.com/Shubiks/go-simple-api/internal/config"
	"github.com/Shubiks/go-simple-api/internal/db"
	"github.com/Shubiks/go-simple-api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/Shubiks/go-simple-api/internal/s3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system env")
	}
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	if err := s3.InitS3(); err != nil {
		log.Fatalf("Failed to initialize AWS S3: %v", err)
	}

	handler.SetDB(database)
	handler.SetFollowDB(database.DB)

	r := chi.NewRouter()
	r.Use(middleware.Logger) // Optional: logs requests

	// 🔓 Public routes
	r.Group(func(r chi.Router) {
		r.Post("/users", handler.CreateUserHandler)
		r.Post("/login", handler.LoginHandler)
		r.Post("/refresh", handler.RefreshHandler)
	})

	// 🔐 Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware) // ✅ Apply JWT middleware only here
		r.Get("/getusers", handler.GetUsersHandler(database))
		r.Post("/follow/{user_id}", handler.SendFollowRequestHandler)
		r.Post("/upload/profile-picture", handler.UploadProfilePictureHandler)
	})

	log.Printf("Server starting on port %s...", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
