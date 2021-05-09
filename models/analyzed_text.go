package models

import (
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/nlp"
	"time"
)

type AnalyzedText struct {
	Content  string    `json:"content" validate:"required"`
	Author   string    `json:"author" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
	People   []string  `json:"people" validate:"required"`
	GPE      []string  `json:"gpe" validate:"required"`
}

func NewAnalyzedText(text Text) (analyzedText AnalyzedText, err error) {
	var people, gpe []string
	entities := nlp.ExtractEntities(text.Content)

	for _, entity := range entities {
		switch entity.Type {
		case "PERSON":
			people = append(people, entity.Text)
		case "GPE":
			gpe = append(gpe, entity.Text)
		}
	}

	analyzedText = AnalyzedText{
		Content:  text.Content,
		Author:   text.Author,
		Date:     text.Date,
		SourceID: text.SourceID,
		People:   people,
		GPE:      gpe,
	}

	_, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT"), analyzedText, text.ID)
	if err == nil {
		text.Analyzed = true
		elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, text.ID)
	}

	return analyzedText, err
}
