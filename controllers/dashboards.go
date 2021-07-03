package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type PostDashboardParams struct {
	Title    string `json:"title" binding:"required"`
	SourceID int    `json:"sourceID" binding:"required"`
}

func PostDashboard(context *gin.Context) {
	var params PostDashboardParams
	var source models.Source

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

	dashboard, err := models.NewDashboard(params.Title, int(source.ID))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":      "Dashboard created",
		"dashboardID": dashboard.ID,
	})
}

type GetDashboardsParams struct {
	SourceID int `form:"sourceID" binding:"required"`
}

func GetDashboards(context *gin.Context) {
	var params GetDashboardsParams
	var source models.Source
	var dashboards []models.Dashboard
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

	database.Database.Where("source_id = ?", source.ID).Order("updated_at desc").Order("updated_at desc").Find(&dashboards)

	for _, dashboard := range dashboards {
		result = append(result, map[string]interface{}{
			"last_updated": dashboard.UpdatedAt.Format(insightDateFormat),
			"title":        dashboard.Title,
			"id":           dashboard.ID,
			"token":        dashboard.Token,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"dashboards": result,
	})
}

type DeleteDashboardParams struct {
	SourceID    int `form:"sourceID" binding:"required"`
	DashboardID int `form:"dashboardID" binding:"required"`
}

func DeleteDashboard(context *gin.Context) {
	var params DeleteDashboardParams
	var source models.Source
	var dashboard models.Dashboard

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

	lookupResult := database.Database.Where("source_id = ? and id = ?", source.ID, params.DashboardID).Find(&dashboard)
	if lookupResult.Error == nil {
		database.Database.Delete(&dashboard)
		context.JSON(http.StatusOK, gin.H{
			"status": "dashboard deleted",
		})
	} else {
		context.JSON(http.StatusPreconditionFailed, gin.H{
			"status": "dashboard could not be deleted",
		})
	}
}

type AddInsightToDashboardParams struct {
	InsightID   int `json:"insightID" binding:"required"`
	DashboardID int `json:"dashboardID" binding:"required"`
	SourceID    int `json:"sourceID" binding:"required"`
}

func AddInsightToDashboard(context *gin.Context) {
	var params AddInsightToDashboardParams
	var source models.Source
	var insight models.Insight
	var dashboard models.Dashboard

	user, _ := context.Get("user")
	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, params.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	database.Database.Where("source_id = ? and id = ?", source.ID, params.DashboardID).Find(&dashboard)
	if dashboard.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Dashboard could not be validated"})
		return
	}

	database.Database.Where("source_id = ? and id = ?", source.ID, params.InsightID).Find(&insight)
	if insight.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Insight could not be validated"})
		return
	}

	_, err := models.NewDashboardInsight(int(dashboard.ID), int(insight.ID))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Insight added to dashboard",
	})
}
