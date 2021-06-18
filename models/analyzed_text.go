package models

import (
	"log"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/nlp"
	"time"
)

type AnalyzedText struct {
	ID       string    `json:"-"`
	Content  string    `json:"content" validate:"required"`
	Author   []string  `json:"author" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
	People   []string  `json:"people" validate:"required"`
	GPE      []string  `json:"gpe" validate:"required"`
	Tokens   []string  `json:"tokens" validate:"required"`
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

	tokens, err := nlp.Tokenize(text.Content)
	if err != nil {
		return analyzedText, err
	}

	analyzedText = AnalyzedText{
		Content:  text.Content,
		Author:   text.Author,
		Date:     text.Date,
		SourceID: text.SourceID,
		People:   people,
		GPE:      gpe,
		Tokens:   tokens,
	}

	_, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT"), analyzedText, text.ID)
	if err == nil {
		text.Analyzed = true
		elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, text.ID)
	}

	return analyzedText, err
}

func GetAnalyzedTexts(
	searchText string, searchAuthor []string, people []string, gpe []string,
	startDate time.Time, endDate time.Time, sourceID int,
) (analyzedTexts []AnalyzedText, err error) {
	analyzedTexts = []AnalyzedText{}

	textQuery := elastic.NewAnalyzedTextQuery(searchText, searchAuthor, people, gpe, startDate, endDate, sourceID)
	if queryResult, err := elastic.Query(textQuery, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return analyzedTexts, err
	} else {
		for _, hit := range queryResult.Hits.Hits {
			analyzedTexts = append(analyzedTexts, AnalyzedText{
				ID:       hit.ID,
				Content:  hit.Source.Content,
				Author:   hit.Source.Author,
				SourceID: hit.Source.SourceID,
				People:   hit.Source.People,
				GPE:      hit.Source.GPE,
				Tokens:   hit.Source.Tokens,
				// Date:     hit.Source.Date,
			})
		}
	}

	return analyzedTexts, err
}
