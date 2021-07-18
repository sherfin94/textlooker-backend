package models

import (
	"encoding/json"
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

func CreateGeneralAggregationFromQueryResult(queryResult elastic.QueryResult) (aggregation map[string]interface{}) {
	var queryResultMap map[string]interface{}
	bytes, err := json.Marshal(queryResult.AggregationsPart)
	if err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(bytes, &queryResultMap)
	}

	aggregation = map[string]interface{}{}
	for _, field := range elastic.AggregatableFields {
		buckets := queryResultMap[field].(map[string]interface{})["buckets"].([]interface{})
		countSet := []map[string]interface{}{}
		for _, bucket := range buckets {
			count := map[string]interface{}{
				"key":   bucket.(map[string]interface{})["key"],
				"count": bucket.(map[string]interface{})["doc_count"],
			}
			countSet = append(countSet, count)
		}
		aggregation[field] = countSet
	}

	return aggregation
}

func CreatePerDateAggregationFromQueryResult(queryResult elastic.QueryResult, field string) (counts []CountItem) {
	counts = []CountItem{}
	// log.Println(queryResult.AggregationsPart.(map[string]interface{})["per_date"].(map[string]interface{})["buckets"].([]map[string]interface{}))
	for _, bucket := range queryResult.AggregationsPart.(map[string]interface{})["per_date"].(map[string]interface{})["buckets"].([]interface{}) {
		key := bucket.(map[string]interface{})["key"].(map[string]interface{})
		counts = append(counts, CountItem{
			Date:  util.ParseTimestamp(key["date"].(float64)),
			Value: key[field].(string),
			Count: int(bucket.(map[string]interface{})["doc_count"].(float64)),
		})
	}

	return counts
}

func GetAggregation(
	searchText string, filterItems []elastic.FilterItem,
	startDate time.Time, endDate time.Time, sourceID int, dateAvailableForSource bool,
) (aggregation map[string]interface{}, err error) {

	query := elastic.NewAggregateAllQuery(
		searchText, filterItems, startDate,
		endDate, sourceID, dateAvailableForSource,
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
	searchText string, filterItems []elastic.FilterItem,
	startDate time.Time, endDate time.Time, sourceID int, field string, dateAvailableForSource bool,
) (counts []CountItem, err error) {

	query := elastic.NewAggregateByOneFieldQuery(
		searchText, filterItems, field, startDate,
		endDate, sourceID, dateAvailableForSource,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return counts, err
	} else {
		counts = CreatePerDateAggregationFromQueryResult(queryResult, field)
	}

	return counts, err
}

func GetDatelessAggregation(
	searchText string, filterItems []elastic.FilterItem, sourceID int, dateAvailableForSource bool,
) (aggregation map[string]interface{}, err error) {

	query := elastic.NewDatelessAggregateAllQuery(
		searchText, filterItems, sourceID, dateAvailableForSource,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		return aggregation, err
	} else {
		aggregation = CreateGeneralAggregationFromQueryResult(queryResult)
	}

	return aggregation, err
}
