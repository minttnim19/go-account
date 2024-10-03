package main

import (
	"go-account/config"
	"go-account/internal/api"
	"go-account/internal/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.InitMongo(cfg.MongoDBURI, cfg.MongoDBUser, cfg.MongoDBPassword, cfg.MongoDBName)

	r := api.SetupRouter(db) // Setup routes
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
