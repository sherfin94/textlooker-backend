package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/elastic"
	"textlooker-backend/handlers"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const ReferenceDate = "2006-01-02T15:04:05-07:00"

type BatchTextParams struct {
	Batch    []TextParams `json:"batch"`
	SourceID int          `json:"sourceID" validate:"required"`
}

type TextParams struct {
	Content string   `json:"content" validate:"required"`
	Author  []string `json:"author,omitempty" validate:"required"`
	Date    string   `json:"date,omitempty" validate:"required"`
}

func PostText(context *gin.Context) {
	var source models.Source
	var batchParams BatchTextParams
	lastOccuredError := ""
	user, _ := context.Get("user")

	if err := context.ShouldBindJSON(&batchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, batchParams.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	var textBatch handlers.TextBatch
	for _, textParams := range batchParams.Batch {
		text := handlers.Text{
			Content:      textParams.Content,
			Author:       textParams.Author,
			DateAsString: textParams.Date,
		}
		textBatch.TextSet = append(textBatch.TextSet, text)
	}

	textBatch.SourceID = batchParams.SourceID
	count, err := handlers.ProcessTextBatch(textBatch, &source)

	if err != nil {
		lastOccuredError = err.Error()
	}

	context.JSON(http.StatusOK, gin.H{
		"savedTextCount":   count,
		"lastOccuredError": lastOccuredError,
	})
}

type TextSearchParams struct {
	Content   string       `form:"content,default=*"`
	Filter    []FilterItem `form:"filter"`
	StartDate string       `form:"startDate" validate:"required"`
	EndDate   string       `form:"endDate" validate:"required"`
	SourceID  int          `form:"sourceID" validate:"required"`
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

	var filterItems []elastic.FilterItem
	for _, item := range textSearchParams.Filter {
		filterItems = append(filterItems, elastic.FilterItem{Label: item.Label, Text: item.Text})
	}
	texts, err := models.GetTexts(
		textSearchParams.Content,
		filterItems,
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
