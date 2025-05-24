package main

import (
	"log"
	"net/http"

	"github.com/Shubiks/go-simple-api/internal/auth"
	"github.com/Shubiks/go-simple-api/internal/config"
	"github.com/Shubiks/go-simple-api/internal/db"
	"github.com/Shubiks/go-simple-api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// 👇 load .env file first
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system env")
	}
	cfg := config.Load()

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	handler.SetDB(db)

	r := chi.NewRouter()
	r.Use(auth.JWTMiddleware)

	r.Post("/users", handler.CreateUserHandler)
	r.Post("/login", handler.LoginHandler)

	r.Get("/users", handler.GetUsersHandler(db))

	log.Printf("Server starting on port %s...", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
