package controllers

import (
	"net/http"
	"textlooker-backend/elastic"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type FilterItem struct {
	Label string `form:"label" validate:"required"`
	Text  string `form:"text" validate:"required"`
}

type AnalyzedTextSearchParams struct {
	Content   string       `form:"content,default=*"`
	StartDate string       `form:"startDate" validate:"required"`
	EndDate   string       `form:"endDate" validate:"required"`
	SourceID  int          `form:"sourceID" validate:"required"`
	Filter    []FilterItem `form:"filter[]"`
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

	var filterItems []elastic.FilterItem
	for _, item := range analyzedTextSearchParams.Filter {
		filterItems = append(filterItems, elastic.FilterItem{Label: item.Label, Text: item.Text})
	}
	texts, err := models.GetAnalyzedTexts(
		analyzedTextSearchParams.Content,
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
