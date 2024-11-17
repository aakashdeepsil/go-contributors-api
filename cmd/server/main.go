package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aakashdeepsil/go-contributors-api/internal/config"
	"github.com/aakashdeepsil/go-contributors-api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()

	// Initialize repositories
	repos, err := database.NewRepositories(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}
	log.Printf("Repositories initialized: %v", repos)

	// Server setup code...
	server := &http.Server{
		Addr: ":" + cfg.Port,
		// Handler will be set later
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
