package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(r *gin.RouterGroup, db *mongo.Database) {

	// API ping
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Pong!"})
	})

	// API Users
	UserRoutes(r, db)
}
