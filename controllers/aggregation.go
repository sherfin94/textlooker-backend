package controllers

import (
	"net/http"
	"textlooker-backend/elastic"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AggregationParams struct {
	AnalyzedTextSearchParams
	Field string `form:"field" validate:"required"`
}

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

func GetPerDateAggregation(context *gin.Context) {
	var params AggregationParams
	var source models.Source
	var startDate, endDate time.Time
	var field elastic.AggregationField

	err := bindAggregationParamsToSourceFieldAndDateRange(context, &params, &source, &startDate, &endDate, &field)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	counts, err := models.GetPerDateAggregation(
		params.Content, params.Author, params.People, params.GPE,
		startDate, endDate, int(source.ID), field,
	)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"aggregation": counts})
}

type DatelessAggregationParams struct {
	Content  string   `form:"content,default=*"`
	Author   []string `form:"author[]"`
	SourceID int      `form:"sourceID" validate:"required"`
	People   []string `form:"people[]"`
	GPE      []string `form:"gpe[]"`
	Tokens   []string `form:"tokens[]"`
}

func GetDatelessGeneralAggregation(context *gin.Context) {
	var params DatelessAggregationParams
	var source models.Source

	err := bindParamsToSource(context, &params, &source)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aggregation, err := models.GetDatelessAggregation(
		params.Content, params.Author, params.People, params.GPE, params.Tokens,
		int(source.ID),
	)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"aggregation": aggregation})
}
