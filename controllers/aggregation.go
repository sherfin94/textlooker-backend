package controllers

import (
	"net/http"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetGeneralAggregation(context *gin.Context) {
	var params AnalyzedTextSearchParams
	var source models.Source
	var startDate, endDate time.Time

	err := bindParamsToSourceAndDateRange(context, &params, &source, &startDate, &endDate)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aggregation, err := models.GetAggregation(
		params.Content, params.Author, params.People, params.GPE,
		startDate, endDate, int(source.ID),
	)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"aggregation": aggregation})
}
