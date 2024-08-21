package main

import (
	"go-account/config"
	"go-account/middlewares"
	"go-account/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Environment variable not set")
	}

	db, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatal("Could not connect to MongoDB")
	}

	r := gin.Default()
	r.Use(gin.Recovery()) // Recovery when system die.
	r.Use(middlewares.ErrorHandler())

	version1 := r.Group("/v1")
	routes.InitRoutes(version1, db)
	r.Run(":8080")
}
