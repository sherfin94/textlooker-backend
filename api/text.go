package api

import (
	"net/http"
	apihandlers "textlooker-backend/api_handlers"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const ReferenceDate = "Jan 2 15:04:05 -0700 MST 2006"

type BatchTextParams struct {
	Batch []TextParams `json:"batch"`
}

type TextParams struct {
	Content string   `json:"content" validate:"required"`
	Author  []string `json:"author,omitempty" validate:"required"`
	Date    string   `json:"date,omitempty" validate:"required"`
}

func PostText(context *gin.Context) {
	var source *models.Source
	var batchParams BatchTextParams

	if err := context.ShouldBindJSON(&batchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sourceData, _ := context.Get("source")
	source = sourceData.(*models.Source)

	count := 0
	lastOccuringErrorMessage := ""
	for _, textParams := range batchParams.Batch {
		date, err := time.Parse("2006-01-02T15:04:05-07:00", textParams.Date)
		if err == nil {
			if err := apihandlers.TextWithDate(
				textParams.Content,
				textParams.Author,
				date,
				source,
			); err == nil {
				count += 1
			} else {
				lastOccuringErrorMessage = err.Error()
			}
		} else {
			if err := apihandlers.TextWithoutDate(
				textParams.Content,
				textParams.Author,
				source,
			); err == nil {
				count += 1
			} else {
				lastOccuringErrorMessage = err.Error()
			}
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"savedTextCount":           count,
		"lastOccuringErrorMessage": lastOccuringErrorMessage,
	})
}
