package main

import (
	"log"

	"github.com/falasefemi2/hms/internal/config"
	"github.com/falasefemi2/hms/internal/database"
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
	log.Println("Application ready!")
}
