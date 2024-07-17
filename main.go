package main

import (
	"go-account/config"
	"go-account/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.Connect2MongoDB()
	if err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
	r := gin.Default()
	r.Use(gin.Recovery()) // Recovery when system die.
	version1 := r.Group("/v1")

	routes.InitRoutes(version1, db)
	r.Run()
}
