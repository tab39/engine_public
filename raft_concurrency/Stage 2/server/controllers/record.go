package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"theantisnipe/ds-stage1-backend/models"

	"github.com/gin-gonic/gin"
)

type createAlbumInput struct {
	Title         string  `json:"title" binding:"required"`
	Artist        string  `json:"artist" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	MasterMachine bool    `json:"is_master"`
}

func FindAlbums(c *gin.Context) {
	var albums []models.Album
	models.DB.Find(&albums)
	c.JSON(http.StatusOK, gin.H{"data": albums})
}

func CreateAlbum(c *gin.Context) {
	input := createAlbumInput{MasterMachine: true}
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.MasterMachine {
		postBody := createAlbumInput{Title: input.Title, Artist: input.Artist, Price: input.Price, MasterMachine: false}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(postBody)
		if err != nil {
			log.Fatal(err)
		}
		_, err = http.Post("http://companion-app-1:8080/albums", "application/json", &buf)
		if err != nil {
			log.Fatal(err)
		}

		err = json.NewEncoder(&buf).Encode(postBody)
		if err != nil {
			log.Fatal(err)
		}

		_, err = http.Post("http://companion-app-2:8080/albums", "application/json", &buf)
		if err != nil {
			log.Fatal(err)
		}
	}
	album := models.Album{Title: input.Title, Artist: input.Artist, Price: input.Price}
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
