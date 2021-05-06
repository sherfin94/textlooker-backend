package controllers

import (
	"net/http"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const referenceTime = "Jan 2 15:04:05 -0700 MST 2006"

type Text struct {
	Content  string `json:"content" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Time     string `json:"time" validate:"required"`
	SourceID int    `json:"sourceID" validate:"required"`
}

func PostText(context *gin.Context) {
	var text Text

	if err := context.ShouldBindJSON(&text); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	time, err := time.Parse(referenceTime, text.Time)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.NewText(text.Content, text.Author, time, text.SourceID); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Text saved",
	})
}
