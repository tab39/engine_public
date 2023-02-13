package main

import (
	"theantisnipe/ds-stage1-backend/controllers"

	"theantisnipe/ds-stage1-backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	router.GET("/albums", controllers.FindAlbums)
	router.GET("/albums/:id", controllers.FindSingleAlbum)
	router.POST("/albums", controllers.CreateAlbum)
	router.Run("localhost:8080")
}
