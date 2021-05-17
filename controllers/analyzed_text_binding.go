package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func bindParamsToSourceAndDateRange(
	context *gin.Context,
	analyzedTextSearchParams *AnalyzedTextSearchParams,
	source *models.Source, startDate *time.Time, endDate *time.Time,
) (err error) {
	if err = context.BindQuery(analyzedTextSearchParams); err != nil {
		return err
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

	(*startDate), err = time.Parse(ReferenceDate, analyzedTextSearchParams.StartDate)
	if err != nil {
		return err
	}
	(*endDate), err = time.Parse(ReferenceDate, analyzedTextSearchParams.EndDate)

	return err
}
