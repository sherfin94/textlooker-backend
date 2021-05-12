package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyzedTextSearchParams struct {
	Content   string   `form:"content,default=*"`
	Author    string   `form:"author,default=*"`
	StartDate string   `form:"startDate" validate:"required"`
	EndDate   string   `form:"endDate" validate:"required"`
	SourceID  int      `form:"sourceID" validate:"required"`
	People    []string `form:"people"`
	GPE       []string `form:"gpe"`
}

func GetAnalyzedTexts(context *gin.Context) {
	var analyzedTextSearchParams AnalyzedTextSearchParams
	var source models.Source

	if err := context.BindQuery(&analyzedTextSearchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where(
		"user_id = ? and id = ?",
		user.(*models.User).ID,
		analyzedTextSearchParams.SourceID).Find(&source)

	if sourceSearchResult.Error != nil || source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be verified`"})
		return
	}

	startDate, err1 := time.Parse(ReferenceDate, analyzedTextSearchParams.StartDate)
	endDate, err2 := time.Parse(ReferenceDate, analyzedTextSearchParams.EndDate)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse either or both of the dates"})
		return
	}

	texts, err := models.GetAnalyzedTexts(
		analyzedTextSearchParams.Content,
		analyzedTextSearchParams.Author,
		analyzedTextSearchParams.People,
		analyzedTextSearchParams.GPE,
		startDate, endDate,
		analyzedTextSearchParams.SourceID,
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"texts": texts})
	}
}
