package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type GetInsightsOfDashboardParams struct {
	DashboardID int    `form:"dashboardID"`
	Token       string `form:"token"`
}

func GetInsightsOfDashboardViaToken(context *gin.Context) {
	var params GetInsightsOfDashboardParams
	var dashboard models.Dashboard
	var dashboard_insights []models.DashboardInsight
	var result []int

	if err := context.BindQuery(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lookupResult := database.Database.Where("id = ?", dashboard.ID).Find(&dashboard)
	if lookupResult.Error == nil {

		if dashboard.Token != params.Token {
			context.JSON(http.StatusPreconditionFailed, gin.H{
				"status": "Invalid token",
			})
			return
		}

		database.Database.Where("dashboard_id = ?", dashboard.ID).Order("updated_at desc").Order("updated_at desc").Find(&dashboard_insights)

		for _, dashboard_insight := range dashboard_insights {
			result = append(result, int(dashboard_insight.InsightID))
		}

		context.JSON(http.StatusOK, gin.H{
			"insightIDs": result,
		})
	} else {
		context.JSON(http.StatusPreconditionFailed, gin.H{
			"status": "insight IDs could not be fetched",
		})
	}
}

type GetInsightsAggregationViaTokenParams struct {
	DashboardID int    `form:"dashboardID"`
	InsightID   int    `form:"insightID"`
	Token       string `form:"token"`
}

func GetInsightsAggregationViaToken(context *gin.Context) {
	var params GetInsightsAggregationViaTokenParams
	var dashboard models.Dashboard
	var insight models.Insight
	var dashboard_insight models.DashboardInsight

	if err := context.BindQuery(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lookupResult := database.Database.Where("id = ?", params.DashboardID).Find(&dashboard)
	if dashboard.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": lookupResult.Error.Error()})
		return
	}

	if dashboard.Token != params.Token {
		context.JSON(http.StatusPreconditionFailed, gin.H{
			"status": "Invalid token",
		})
		return
	}

	lookupResult = database.Database.Where("insight_id = ? and dashboard_id = ?", params.InsightID, dashboard.ID).Find(&dashboard_insight)
	if dashboard_insight.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": lookupResult.Error.Error()})
		return
	}

	lookupResult = database.Database.Where("id = ?", dashboard_insight.InsightID).Find(&insight)
	if dashboard_insight.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": lookupResult.Error.Error()})
		return
	}

	aggregation, err := insight.Aggregation()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"aggregation": aggregation,
	})
}
