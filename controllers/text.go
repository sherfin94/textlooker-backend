package controllers

import (
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/handlers"
	"textlooker-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

const ReferenceDate = "Jan 2 15:04:05 -0700 MST 2006"

type BatchTextParams struct {
	Batch    []TextParams `json:"batch"`
	SourceID int          `json:"sourceID" validate:"required"`
}

type TextParams struct {
	Content string   `json:"content" validate:"required"`
	Author  []string `json:"author,omitempty" validate:"required"`
	Date    string   `json:"date,omitempty" validate:"required"`
}

func PostText(context *gin.Context) {
	var source models.Source
	var batchParams BatchTextParams
	user, _ := context.Get("user")

	if err := context.ShouldBindJSON(&batchParams); err != nil {
		println("shashi")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, batchParams.SourceID).Find(&source)
	if source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be validated"})
		return
	}

	count := 0
	for _, textParams := range batchParams.Batch {
		if err := handlers.Text(
			textParams.Content,
			textParams.Author,
			textParams.Date,
			int(source.ID),
		); err == nil {
			count += 1
		} else {
			println(err.Error())
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"savedTextCount": count,
	})
}

type TextSearchParams struct {
	Content   string   `form:"content,default=*"`
	Author    []string `form:"author[],default="`
	StartDate string   `form:"startDate" validate:"required"`
	EndDate   string   `form:"endDate" validate:"required"`
	SourceID  int      `form:"sourceID" validate:"required"`
}

func GetTexts(context *gin.Context) {
	var textSearchParams TextSearchParams
	var source models.Source

	if err := context.BindQuery(&textSearchParams); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	sourceSearchResult := database.Database.Where("user_id = ? and id = ?", user.(*models.User).ID, textSearchParams.SourceID).Find(&source)
	if sourceSearchResult.Error != nil || source.ID == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Source could not be verified`"})
		return
	}

	startDate, err1 := time.Parse(ReferenceDate, textSearchParams.StartDate)
	endDate, err2 := time.Parse(ReferenceDate, textSearchParams.EndDate)

	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse either or both of the dates"})
		return
	}

	texts, err := models.GetTexts(
		textSearchParams.Content,
		textSearchParams.Author,
		startDate, endDate,
		int(source.ID),
	)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"texts": texts})
	}
}
