package models

import (
	"log"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/util"
	"time"
)

type CountItem struct {
	Date  *time.Time  `json:"date,omitempty" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
	Count int         `json:"count" validate:"required"`
}

type Aggregation struct {
	Authors []CountItem `json:"authors" validate:"required"`
	People  []CountItem `json:"people" validate:"required"`
	GPE     []CountItem `json:"gpe" validate:"required"`
	Tokens  []CountItem `json:"tokens" validate:"required"`
	Dates   []CountItem `json:"dates,omitempty" validate:"required"`
}

func CreateGeneralAggregationFromQueryResult(queryResult elastic.QueryResult) (aggregation Aggregation) {
	authors, people, gpe, tokens, dates := []CountItem{}, []CountItem{}, []CountItem{}, []CountItem{}, []CountItem{}

	for _, bucket := range queryResult.AggregationsPart.AuthorAggregation.Buckets {
		authors = append(authors, CountItem{Value: bucket.Key, Count: bucket.Value})
	}

	for _, bucket := range queryResult.AggregationsPart.PeopleAggregation.Buckets {
		people = append(people, CountItem{Value: bucket.Key, Count: bucket.Value})
	}

	for _, bucket := range queryResult.AggregationsPart.GPEAggregation.Buckets {
		gpe = append(gpe, CountItem{Value: bucket.Key, Count: bucket.Value})
	}

	for _, bucket := range queryResult.AggregationsPart.TokenAggregation.Buckets {
		tokens = append(tokens, CountItem{Value: bucket.Key, Count: bucket.Value})
	}

	for _, bucket := range queryResult.AggregationsPart.DateAggregation.Buckets {
		dates = append(dates, CountItem{Value: util.ParseTimestamp(bucket.Key.(float64)), Count: bucket.Value})
	}

	return Aggregation{
		Authors: authors,
		People:  people,
		GPE:     gpe,
		Tokens:  tokens,
		Dates:   dates,
	}
}

func CreatePerDateAggregationFromQueryResult(queryResult elastic.QueryResult) (counts []CountItem) {
	counts = []CountItem{}
	for _, bucket := range queryResult.AggregationsPart.PerDateAggregation.Buckets {
		counts = append(counts, CountItem{Value: bucket.Key.FieldValue, Count: bucket.Count, Date: util.ParseTimestamp(bucket.Key.Date)})
	}

	return counts
}

func GetAggregation(
	searchText string, searchAuthor []string, people []string, gpe []string,
	startDate time.Time, endDate time.Time, sourceID int,
) (aggregation Aggregation, err error) {

	query := elastic.NewAggregateAllQuery(
		searchText, searchAuthor, people, gpe, startDate,
		endDate, sourceID,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return aggregation, err
	} else {
		aggregation = CreateGeneralAggregationFromQueryResult(queryResult)
	}

	return aggregation, err
}

func GetPerDateAggregation(
	searchText string, searchAuthor []string, people []string, gpe []string,
	startDate time.Time, endDate time.Time, sourceID int, field elastic.AggregationField,
) (counts []CountItem, err error) {

	query := elastic.NewAggregateByOneFieldQuery(
		searchText, searchAuthor, people, gpe, startDate,
		endDate, sourceID, field,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return counts, err
	} else {
		counts = CreatePerDateAggregationFromQueryResult(queryResult)
	}

	return counts, err
}

func GetDatelessAggregation(
	searchText string, searchAuthor []string, people []string, gpe []string, sourceID int,
) (aggregation Aggregation, err error) {

	query := elastic.NewDatelessAggregateAllQuery(
		searchText, searchAuthor, people, gpe, sourceID,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return aggregation, err
	} else {
		aggregation = CreateGeneralAggregationFromQueryResult(queryResult)
	}

	return aggregation, err
}
