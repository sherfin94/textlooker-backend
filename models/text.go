package models

import (
	"log"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"time"

	"github.com/go-playground/validator/v10"
)

type Text struct {
	ID       string    `json:"-"`
	Content  string    `json:"content" validate:"required"`
	Author   []string  `json:"author" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
	Analyzed bool      `json:"analyzed"`
}

func NewText(content string, author []string, date time.Time, sourceID int) (text Text, err error) {
	text = Text{Content: content, Author: author, Date: date, SourceID: sourceID, Analyzed: false}
	validator := validator.New()
	if err = validator.Struct(text); err != nil {
		return text, err
	}

	if text.ID, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, ""); err != nil {
		return text, err
	}

	return text, nil
}

func GetTexts(content string, author []string, dateStart time.Time, dateEnd time.Time, sourceID int) (texts []Text, err error) {
	textQuery := elastic.NewTextQuery(content, author, dateStart, dateEnd, sourceID)
	texts = []Text{}

	if err != nil {
		log.Println(err)
		return texts, err
	}
	if queryResult, err := elastic.Query(textQuery, deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT")); err != nil {
		log.Println(err)
		return texts, err
	} else {
		for _, hit := range queryResult.Hits.Hits {
			texts = append(texts, Text{
				ID:       hit.ID,
				Content:  hit.Source.Content,
				Author:   hit.Source.Author,
				SourceID: hit.Source.SourceID,
				Analyzed: hit.Source.Analyzed,
				Date:     hit.Source.Date,
			})
		}
	}
	return texts, err
}
