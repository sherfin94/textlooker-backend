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
	}

	_, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT"), analyzedText, text.ID)
	if err == nil {
		text.Analyzed = true
		elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, text.ID)
	}

	return analyzedText, err
}

func GetAnalyzedTexts(
	searchText string, from int, filterItems []elastic.FilterItem,
	startDate time.Time, endDate time.Time, sourceID int,
	dateRangeProvided bool,
) (analyzedTexts []AnalyzedText, total int, err error) {
	analyzedTexts = []AnalyzedText{}

	textQuery := elastic.NewAnalyzedTextQuery(searchText, filterItems, startDate, endDate, sourceID, dateRangeProvided)
	textQuery.Size = 20
	textQuery.From = from

	log.Println(textQuery.RequestString())

	total = 0
	if queryResult, err := elastic.Query(textQuery, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return analyzedTexts, total, err
	} else {
		for _, hit := range queryResult.Hits.Hits {

			// date, err := time.Parse("2006-01-02T15:04:05-0700", hit.Source.Date)
			// if err != nil {
			// 	return analyzedTexts, err
			// }

			// log.Println("shashi")
			// log.Println(date.String())
			analyzedTexts = append(analyzedTexts, AnalyzedText{
				ID:       hit.ID,
				Content:  hit.Source.Content,
				Author:   hit.Source.Author,
				SourceID: hit.Source.SourceID,
				Date:     hit.Source.Date.Time,
				// People:   hit.Source.People,
				// GPE:      hit.Source.GPE,
				// Tokens:   hit.Source.Tokens,
			})
		}
		total = queryResult.Hits.Total.Value
	}

	return analyzedTexts, total, err
}
