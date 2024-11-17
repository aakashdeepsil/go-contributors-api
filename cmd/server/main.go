package main

import (
	"log"
	"net/http"

	"github.com/aakashdeepsil/go-contributors-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	server := &http.Server{
		Addr: ":" + cfg.Port,
		// Handler will be set later
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
