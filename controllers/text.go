package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const ReferenceDate = "Jan 2 15:04:05 -0700 MST 2006"

type Text struct {
	Content  string   `json:"content" validate:"required"`
	Author   []string `json:"author" validate:"required"`
	Date     string   `json:"date" validate:"required"`
	SourceID int      `json:"sourceID" validate:"required"`
}

func PostText(context *gin.Context) {
	var text Text
	var source models.Source

	if err := context.ShouldBindJSON(&text); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	time, err := time.Parse(ReferenceDate, text.Date)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, text.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source not found"})
		return
	}

	if text, err := models.NewText(text.Content, text.Author, time, int(source.ID)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		go models.NewAnalyzedText(text)
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Text saved",
	})
}

type TextSearchParams struct {
	Content   string   `form:"content,default=*"`
	Author    []string `form:"author[],default="`
	StartDate string   `form:"startDate" validate:"required"`
	EndDate   string   `form:"endDate" validate:"required"`
	SourceID  int      `form:"sourceID" validate:"required"`
}

func GetTexts(context *gin.Context) {
	var textSearchParams TextSearchParams
	var source models.Source

	if err := context.BindQuery(&textSearchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, textSearchParams.SourceID).Find(&source)
	if sourceSearchResult.Error != nil || source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be verified`"})
		return
	}

	startDate, err1 := time.Parse(ReferenceDate, textSearchParams.StartDate)
	endDate, err2 := time.Parse(ReferenceDate, textSearchParams.EndDate)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse either or both of the dates"})
		return
	}

	texts, err := models.GetTexts(
		textSearchParams.Content,
		textSearchParams.Author,
		startDate, endDate,
		int(source.ID),
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"texts": texts})
	}
}
