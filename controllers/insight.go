package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const insightDateFormat = "2006-01-02T15:04:05-07:00"

type InsightParams struct {
	Title              string `json:"title" binding:"required"`
	SourceID           int    `json:"sourceID" binding:"required"`
	Filter             string `json:"filter" binding:"required"`
	LookForHandle      string `json:"lookForHandle" binding:"required"`
	VisualizeTexts     string `json:"visualizeTexts" binding:"required"`
	VisualizationType  string `json:"visualizationType" binding:"required"`
	StartDate          string `json:"startDate" binding:"required"`
	EndDate            string `json:"endDate" binding:"required"`
	DateRangeAvailable bool   `json:"dateRangeAvailable"`
}

func PostInsight(context *gin.Context) {
	var source models.Source
	var params InsightParams
	var startDate, endDate time.Time

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

	if params.DateRangeAvailable {
		var err error
		startDate, err = time.Parse(insightDateFormat, params.StartDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		endDate, err = time.Parse(insightDateFormat, params.EndDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		startDate = time.Now()
		endDate = time.Now()
	}

	if insight, err := models.NewInsight(
		params.Title,
		params.Filter,
		params.LookForHandle,
		params.VisualizeTexts,
		startDate,
		endDate,
		params.VisualizationType,
		params.DateRangeAvailable,
		int(source.ID),
	); err != nil {
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
	SourceID int `form:"sourceID" binding:"required"`
}

func GetInsights(context *gin.Context) {
	var params GetInsightsParams
	var source models.Source
	var insights []models.Insight
	result := []map[string]interface{}{}

	if err := context.BindQuery(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, params.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	database.Database.Where("source_id = ?", source.ID).Order("updated_at desc").Order("updated_at desc").Find(&insights)

	for _, insight := range insights {
		result = append(result, map[string]interface{}{
			"last_updated":       insight.UpdatedAt.Format(insightDateFormat),
			"title":              insight.Title,
			"id":                 insight.ID,
			"filter":             insight.Filter,
			"lookForHandle":      insight.LookForHandle,
			"visualizeTexts":     insight.VisualizeTexts,
			"startDate":          insight.StartDate.Format(insightDateFormat),
			"endDate":            insight.EndDate.Format(insightDateFormat),
			"dateRangeAvailable": insight.DateRangeAvailable,
			"visualizationType":  insight.VisualizationType,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"insights": result,
	})
}

type DeleteInsightsParams struct {
	SourceID  int `form:"sourceID" binding:"required"`
	InsightID int `form:"insightID" binding:"required"`
}

func DeleteInsight(context *gin.Context) {
	var params DeleteInsightsParams
	var source models.Source
	var insight models.Insight

	if err := context.BindQuery(&params); err != nil {
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
