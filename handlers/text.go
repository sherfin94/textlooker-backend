package handlers

import (
	"strconv"
	"textlooker-backend/models"
	"textlooker-backend/util"
	"time"
)

type TextBatch struct {
	TextSet  []Text `json:"batch"`
	SourceID int    `json:"sourceID" validate:"required"`
}

type Text struct {
	Content      string   `json:"content" validate:"required"`
	Author       []string `json:"author,omitempty" validate:"required"`
	DateAsString string   `json:"-"`
}

func ProcessTextBatch(batch TextBatch, source *models.Source) (int, error) {
	var lastOccuredError error
	var textSet []models.Text
	var date time.Time
	for _, handlerText := range batch.TextSet {
		if source.DateAvailable {
			dateAsInteger, err := strconv.ParseInt(handlerText.DateAsString, 10, 64)
			if err == nil {
				date = *util.ParseTimestamp(float64(dateAsInteger))
			} else {
				lastOccuredError = err
			}
		} else {
			date = time.Now()

		}
		text := models.Text{
			Content:   handlerText.Content,
			Author:    handlerText.Author,
			Date:      date,
			SourceID:  int(source.ID),
			Analyzed:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		textSet = append(textSet, text)
	}

	count, err := models.BulkSaveText(textSet)
	if err != nil {
		return count, err
	} else {
		models.SendToProcessQueue(textSet)
	}

	return count, lastOccuredError
}
