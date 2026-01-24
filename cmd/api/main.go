package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/falasefemi2/hms/docs"
	"github.com/falasefemi2/hms/internal/config"
	"github.com/falasefemi2/hms/internal/database"
	"github.com/falasefemi2/hms/internal/server"
)

// @title Hospital Management System API
// @version 1.0.0
// @description A comprehensive Hospital Management System with user authentication, doctor management, and patient care features
// @termsOfService http://example.com/terms/
// @contact.name HMS Support
// @contact.url http://example.com/support
// @contact.email support@hms.example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @schemes http https
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIs..."
func main() {
	// Load config
	cfg := config.LoadConfig()
	log.Println("Configuration loaded")

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := database.NewDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()
	log.Println("Database connected!")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.InitializeSchema(ctx); err != nil {
		log.Fatalf("Schema initialization error: %v", err)
	}
	log.Println("Database schema initialized!")

	srv := server.NewServer(db)
	port := fmt.Sprintf("%d", cfg.ServerPort)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)

	go func() {
		<-shutdown
		log.Println("Received shutdown signal, shutting down gracefully...")
		if err := srv.Shutdown(15 * time.Second); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("Starting server on port %s...", port)
	if err := srv.Start(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
