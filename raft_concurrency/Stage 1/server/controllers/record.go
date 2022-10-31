package controllers

import (
	"net/http"

	"theantisnipe/ds-stage1-backend/models"

	"github.com/gin-gonic/gin"
)

type createAlbumInput struct {
	ID     int16   `json:"id" binding:"required"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

func FindAlbums(c *gin.Context) {
	var albums []models.Album
	models.DB.Find(&albums)
	c.JSON(http.StatusOK, gin.H{"data": albums})
}

func CreateAlbum(c *gin.Context) {
	var input createAlbumInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	album := models.Album{ID: input.ID, Title: input.Title, Artist: input.Artist, Price: input.Price}
	models.DB.Create(&album)
	c.JSON(http.StatusOK, gin.H{"data": album})
}

func FindSingleAlbum(c *gin.Context) {
	var album models.Album

	if err := models.DB.Where("id = ?", c.Param("id")).First(&album).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": album})
}
