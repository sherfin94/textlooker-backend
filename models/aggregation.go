package models

import (
	"encoding/json"
	"log"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
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

	// for _, bucket := range queryResult.AggregationsPart.AuthorAggregation.Buckets {
	// 	authors = append(authors, CountItem{Value: bucket.Key, Count: bucket.Value})
	// }

	// for _, bucket := range queryResult.AggregationsPart.PeopleAggregation.Buckets {
	// 	people = append(people, CountItem{Value: bucket.Key, Count: bucket.Value})
	// }

	// for _, bucket := range queryResult.AggregationsPart.GPEAggregation.Buckets {
	// 	gpe = append(gpe, CountItem{Value: bucket.Key, Count: bucket.Value})
	// }

	// for _, bucket := range queryResult.AggregationsPart.TokenAggregation.Buckets {
	// 	tokens = append(tokens, CountItem{Value: bucket.Key, Count: bucket.Value})
	// }

	// for _, bucket := range queryResult.AggregationsPart["PERSON"].(AggregationResultPart) {
	// 	dates = append(dates, CountItem{Value: util.ParseTimestamp(bucket.Key.(float64)), Count: bucket.Value})
	// }

	// for _, bucket := range queryResult.AggregationsPart.DateAggregation.Buckets {
	// 	dates = append(dates, CountItem{Value: util.ParseTimestamp(bucket.Key.(float64)), Count: bucket.Value})
	// }

	return aggregation
}

func CreatePerDateAggregationFromQueryResult(queryResult elastic.QueryResult) (counts []CountItem) {
	// counts = []CountItem{}
	// for _, bucket := range queryResult.AggregationsPart.PerDateAggregation.Buckets {
	// 	counts = append(counts, CountItem{Value: bucket.Key.FieldValue, Count: bucket.Count, Date: util.ParseTimestamp(bucket.Key.Date)})
	// }

	return counts
}

func GetAggregation(
	searchText string, filterItems []elastic.FilterItem,
	startDate time.Time, endDate time.Time, sourceID int,
) (aggregation map[string]interface{}, err error) {

	query := elastic.NewAggregateAllQuery(
		searchText, filterItems, startDate,
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
	searchText string, filterItems []elastic.FilterItem,
	startDate time.Time, endDate time.Time, sourceID int, field string,
) (counts []CountItem, err error) {

	query := elastic.NewAggregateByOneFieldQuery(
		searchText, filterItems, field, startDate,
		endDate, sourceID,
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
	searchText string, filterItems []elastic.FilterItem, sourceID int,
) (aggregation map[string]interface{}, err error) {

	query := elastic.NewDatelessAggregateAllQuery(
		searchText, filterItems, sourceID,
	)

	if queryResult, err := elastic.Query(query, deployment.GetEnv("ELASTIC_INDEX_FOR_ANALYZED_TEXT")); err != nil {
		log.Println(err)
		log.Println("ivde ethi")
		return aggregation, err
	} else {
		aggregation = CreateGeneralAggregationFromQueryResult(queryResult)
	}

	return aggregation, err
}
