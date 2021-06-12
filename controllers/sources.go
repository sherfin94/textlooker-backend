package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type Source struct {
	Name            string `json:"name" binding:"required"`
	AuthorAvailable bool   `json:"authorAvailable" binding:"required"`
	DateAvailable   bool   `json:"dateAvailable" binding:"required"`
}

func PostSource(context *gin.Context) {
	var source models.Source

	if err := context.ShouldBindJSON(&source); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")

	if source, err := models.NewSource(source.Name, user.(*models.User), source.DateAvailable, source.AuthorAvailable); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status":   "Source created",
			"sourceID": source.ID,
		})
	}
}

func GetSources(context *gin.Context) {
	var sources []models.Source
	result := []map[string]interface{}{}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ?", user.(*models.User).ID).Order("updated_at desc").Find(&sources)

	for _, source := range sources {
		result = append(result, map[string]interface{}{
			"name":            source.Name,
			"id":              source.ID,
			"authorAvailable": source.AuthorAvailable,
			"dateAvailable":   source.DateAvailable,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"sources": result,
	})
}

func DeleteSource(context *gin.Context) {
	source_id := context.Param("sourceID")
	var source models.Source

	user, _ := context.Get("user")
	result := database.Database.Where("user_id = ? AND id = ?", user.(*models.User).ID, source_id).Find(&source)

	if result.Error == nil {

		database.Database.Delete(&source)
		context.JSON(http.StatusOK, gin.H{
			"status": "source deleted",
		})
	} else {
		context.JSON(http.StatusPreconditionFailed, gin.H{
			"status": "source could not be deleted",
		})
	}
}
