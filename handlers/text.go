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
	Content      string    `json:"content" validate:"required"`
	Author       []string  `json:"author,omitempty" validate:"required"`
	DateAsString string    `json:"-"`
	Date         time.Time `json:"date,omitempty" validate:"required"`
}

func ProcessTextBatch(batch TextBatch, source *models.Source) (int, error) {
	var lastOccuredError error
	for _, text := range batch.TextSet {
		if source.DateAvailable {
			dateAsInteger, err := strconv.ParseInt(text.DateAsString, 10, 64)
			if err == nil {
				date := util.ParseTimestamp(float64(dateAsInteger))
				text.Date = *date
			} else {
				lastOccuredError = err
			}
		} else {
			text.Date = time.Now()
		}
	}

	var textSet []models.Text
	for _, handlerText := range batch.TextSet {
		text := models.Text{
			Content:   handlerText.Content,
			Author:    handlerText.Author,
			Date:      handlerText.Date,
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
		for _, text := range textSet {
			text.SendToProcessQueue()
		}
	}

	return count, lastOccuredError
}
