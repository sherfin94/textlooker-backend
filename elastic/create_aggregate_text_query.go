package elastic

import (
	"time"
)

func NewAggregateAllQuery(content string, filterItems []FilterItem, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	query := NewAnalyzedTextQuery(content, filterItems, startDate, endDate, sourceID)
	query.Size = 0
	query = AddGeneralAggregationPart(query, true)
	return query
}

func NewAggregateByOneFieldQuery(content string, field string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	query := NewAnalyzedTextQuery(content, []FilterItem{}, startDate, endDate, sourceID)
	query.Size = 1000
	query = AddSingleFieldCompositeAggregationPart(query, field)
	return query
}

func NewDatelessAggregateAllQuery(content string, filterItems []FilterItem, sourceID int) TextQuery {
	query := NewDatelessAnalyzedTextQuery(content, filterItems, sourceID)
	query.Size = 0
	query = AddGeneralAggregationPart(query, false)
	return query
}
