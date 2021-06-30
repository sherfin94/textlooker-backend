package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type InsightParams struct {
	Title          string `json:"title" binding:"required"`
	SourceID       int    `json:"sourceID" binding:"required"`
	Filter         string `json:"filter" binding:"required"`
	LookForHandle  string `json:"lookForHandle" binding:"required"`
	VisualizeTexts string `json:"visualizeTexts" binding:"required"`
}

func PostInsight(context *gin.Context) {
	var source models.Source
	var insightParams InsightParams

	if err := context.ShouldBindJSON(&insightParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, insightParams.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	if insight, err := models.NewInsight(insightParams.Title, insightParams.Filter, insightParams.LookForHandle, insightParams.VisualizeTexts, int(source.ID)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status":    "Insight created",
			"insightID": insight.ID,
		})
	}
}

type GetInsightsParams struct {
	SourceID int `json:"sourceID" binding:"required"`
}

func GetInsights(context *gin.Context) {
	var params GetInsightsParams
	var source models.Source
	var insights []models.Insight
	result := []map[string]interface{}{}

	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, params.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	database.Database.Where("source_id = ?", source.ID).Order("updated_at desc").Find(&insights)

	for _, insight := range insights {
		result = append(result, map[string]interface{}{
			"title":          insight.Title,
			"id":             insight.ID,
			"filter":         insight.Filter,
			"lookForHandle":  insight.LookForHandle,
			"visualizeTexts": insight.VisualizeTexts,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"insights": result,
	})
}

type DeleteInsightsParams struct {
	SourceID  int `json:"sourceID" binding:"required"`
	InsightID int `json:"insightID" binding:"required"`
}

func DeleteInsight(context *gin.Context) {
	var params DeleteInsightsParams
	var source models.Source
	var insight models.Insight

	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, params.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	lookupResult := database.Database.Where("source_id = ? and id = ?", source.ID, params.InsightID).Find(&insight)
	if lookupResult.Error == nil {
		database.Database.Delete(&insight)
		context.JSON(http.StatusOK, gin.H{
			"status": "insight deleted",
		})
	} else {
		context.JSON(http.StatusPreconditionFailed, gin.H{
			"status": "insight could not be deleted",
		})
	}
}
