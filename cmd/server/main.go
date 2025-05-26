package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shubiks/go-simple-api/internal/auth"
	"github.com/Shubiks/go-simple-api/internal/config"
	"github.com/Shubiks/go-simple-api/internal/db"
	"github.com/Shubiks/go-simple-api/internal/handler"
	"github.com/Shubiks/go-simple-api/internal/s3"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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
	r.Use(middleware.Logger)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/users", handler.CreateUserHandler)
		r.Post("/login", handler.LoginHandler)
		r.Post("/refresh", handler.RefreshHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		r.Get("/getusers", handler.GetUsersHandler(database))
		r.Post("/follow/{user_id}", handler.SendFollowRequestHandler)
		r.Post("/follow/accept/{user_id}", handler.AcceptFollowRequestHandler)
		r.Post("/upload/profile-picture", handler.UploadProfilePictureHandler)
	})

	// Server setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
