package main

import (
	srv "go-account/cmd/server"
	"go-account/config"
	"go-account/pkg/database"
	"log"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db := database.NewMongoDatabase(cfg)
	srv.NewGonicServer(cfg, db).Start()

}
