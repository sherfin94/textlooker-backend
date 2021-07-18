package handlers

import (
	"errors"
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
	for i := 0; i < 1000 && i < len(batch.TextSet); i++ {
		handlerText := batch.TextSet[i]
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

		authors := []string{}
		for _, authorName := range handlerText.Author {
			authors = append(authors, truncateString(authorName, 1000))
		}

		text := models.Text{
			Content:   truncateString(handlerText.Content, 10000),
			Author:    authors,
			Date:      date,
			SourceID:  int(source.ID),
			Analyzed:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		textSet = append(textSet, text)
	}

	if len(batch.TextSet) > 1000 {
		lastOccuredError = errors.New("batch size greater than 1000. please send data in batches of 1000")
		return 0, lastOccuredError
	}

	count, err := models.BulkSaveText(textSet)
	if err != nil {
		return count, err
	} else {
		models.SendToProcessQueue(textSet)
	}

	return count, lastOccuredError
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		bnoden = str[0:num]
	}
	return bnoden
}
