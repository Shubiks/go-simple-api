package main

import (
	"log"
	"net/http"

	"github.com/Shubiks/go-simple-api/internal/config"
	"github.com/Shubiks/go-simple-api/internal/db"
	"github.com/Shubiks/go-simple-api/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/users", handler.GetUsersHandler(database))

	log.Printf("Server starting on port %s...", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
