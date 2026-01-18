package main

import (
	"fmt"
	"log"

	"github.com/falasefemi2/hms/internal/config"
	"github.com/falasefemi2/hms/internal/database"
	"github.com/falasefemi2/hms/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	log.Println("Configuration loaded")

	err := cfg.Validate()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := database.NewDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	log.Println("Database connected!")

	err = db.InitializeSchema()
	if err != nil {
		log.Fatalf("Schema error: %v", err)
	}

	log.Println("All tables created!")

	// Start the server
	srv := server.NewServer(db)
	err = srv.Start(fmt.Sprintf("%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
