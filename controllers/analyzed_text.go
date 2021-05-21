package controllers

import (
	"net/http"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyzedTextSearchParams struct {
	Content   string   `form:"content,default=*"`
	Author    []string `form:"author[]"`
	StartDate string   `form:"startDate" validate:"required"`
	EndDate   string   `form:"endDate" validate:"required"`
	SourceID  int      `form:"sourceID" validate:"required"`
	People    []string `form:"people[]"`
	GPE       []string `form:"gpe[]"`
}

func GetAnalyzedTexts(context *gin.Context) {
	var analyzedTextSearchParams AnalyzedTextSearchParams
	var source models.Source
	var startDate, endDate time.Time

	err := bindParamsToSourceAndDateRange(context, &analyzedTextSearchParams, &source, &startDate, &endDate)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	texts, err := models.GetAnalyzedTexts(
		analyzedTextSearchParams.Content,
		analyzedTextSearchParams.Author,
		analyzedTextSearchParams.People,
		analyzedTextSearchParams.GPE,
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
