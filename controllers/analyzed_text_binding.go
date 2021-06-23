package controllers

import (
	"errors"
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func bindParamsToSourceAndDateRange(
	context *gin.Context,
	params *AnalyzedTextSearchParams,
	source *models.Source, startDate *time.Time, endDate *time.Time,
) (err error) {
	if err = context.BindQuery(params); err != nil {
		return err
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where(
		"user_id = ? and id = ?",
		user.(*models.User).ID,
		params.SourceID).Find(&source)

	if sourceSearchResult.Error != nil || source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be verified`"})
		return
	}

	(*startDate), err = time.Parse(ReferenceDate, params.StartDate)
	if err != nil {
		return err
	}
	(*endDate), err = time.Parse(ReferenceDate, params.EndDate)

	return err
}

func bindAggregationParamsToSourceFieldAndDateRange(
	context *gin.Context,
	params *AggregationParams,
	source *models.Source, startDate *time.Time, endDate *time.Time,
	field *string,
) (err error) {
	if err = context.BindQuery(params); err != nil {
		return err
	}
	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where(
		"user_id = ? and id = ?",
		user.(*models.User).ID,
		params.SourceID).Find(&source)

	if sourceSearchResult.Error != nil || source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be verified`"})
		return
	}

	(*startDate), err = time.Parse(ReferenceDate, params.StartDate)
	if err != nil {
		return err
	}
	(*endDate), err = time.Parse(ReferenceDate, params.EndDate)
	if err != nil {
		return err
	}

	*field = params.Field

	return err
}

func bindParamsToSource(
	context *gin.Context,
	params *DatelessAggregationParams,
	source *models.Source,
) (err error) {
	if err = context.BindQuery(params); err != nil {
		return err
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where(
		"user_id = ? and id = ?",
		user.(*models.User).ID,
		params.SourceID).Find(&source)

	if sourceSearchResult.Error != nil || source.ID == 0 {
		err = errors.New("Source could not be verified")
		return err
	}

	return err
}
