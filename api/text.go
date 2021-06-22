package api

import (
	"net/http"
	"textlooker-backend/handlers"
	"textlooker-backend/models"

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
	lastOccuredError := ""

	if err := context.ShouldBindJSON(&batchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Request is not formed properly. Please refer to the API documentation."})
		return
	}

	sourceData, _ := context.Get("source")
	source = sourceData.(*models.Source)

	count := 0

	var textBatch handlers.TextBatch
	for _, textParams := range batchParams.Batch {
		text := handlers.Text{
			Content:      textParams.Content,
			Author:       textParams.Author,
			DateAsString: textParams.Date,
		}
		textBatch.TextSet = append(textBatch.TextSet, text)
	}

	textBatch.SourceID = int(source.ID)
	count, err := handlers.ProcessTextBatch(textBatch, source)

	if err != nil {
		lastOccuredError = err.Error()
	}

	context.JSON(http.StatusOK, gin.H{
		"savedTextCount":   count,
		"lastOccuredError": lastOccuredError,
	})
}
