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

	if err := context.ShouldBindJSON(&batchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sourceData, _ := context.Get("source")
	source = sourceData.(*models.Source)

	count := 0
	for _, textParams := range batchParams.Batch {
		if err := handlers.Text(
			textParams.Content,
			textParams.Author,
			textParams.Date,
			source,
		); err == nil {
			count += 1
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"savedTextCount": count,
	})
}
