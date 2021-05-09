package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const referenceDate = "Jan 2 15:04:05 -0700 MST 2006"

type Text struct {
	Content  string `json:"content" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Date     string `json:"date" validate:"required"`
	SourceID int    `json:"sourceID" validate:"required"`
}

func PostText(context *gin.Context) {
	var text Text
	var source Source

	if err := context.ShouldBindJSON(&text); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	time, err := time.Parse(referenceDate, text.Date)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, text.SourceID).Find(&source)
	if sourceSearchResult.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": sourceSearchResult.Error.Error()})
		return
	}

	if text, err := models.NewText(text.Content, text.Author, time, text.SourceID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		go models.NewAnalyzedText(text)
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Text saved",
	})
}
